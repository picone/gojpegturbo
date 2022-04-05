#include "goturbo.h"

// 覆盖原来的format_message方法，因为原来的会打印到控制台。
static void jpeg_err_output_msg(j_common_ptr cinfo) {
    struct my_jpeg_err_mgr* mgr = (struct my_jpeg_err_mgr*)cinfo->err;
    mgr->mgr.format_message(cinfo, mgr->last_msg);
}

// 覆盖原来的error_exit方法，因为原来的错误会调用exit函数导致进程退出。
static void jpeg_err_exit(j_common_ptr cinfo) {
    struct my_jpeg_err_mgr* mgr = (struct my_jpeg_err_mgr*)cinfo->err;
    (*cinfo->err->output_message)(cinfo);
    cinfo->err->num_warnings++;
    longjmp(mgr->setjmp_buf, 1);
}

// 解码jpeg图片
void jpeg_decode(unsigned char* img, unsigned int img_size, jpeg_decode_options* options, jpeg_decode_result* jres) {
    struct jpeg_decompress_struct dinfo;
    my_jpeg_err_mgr               jerr;
    JSAMPROW                      img_decoded = NULL;
    JSAMPROW                      img_row = NULL;
    JSAMPROW                      img_row_start = NULL;
    JSAMPROW                      current_row = NULL;
    JSAMPARRAY                    row_pointer = NULL;
    JDIMENSION                    tmp = 0;
    JDIMENSION                    real_left = 0;
    JDIMENSION                    real_width = 0;
    unsigned int                  crop_width = 0;
    unsigned int                  crop_height = 0;
    size_t                        img_row_size = 0;

    jerr.last_msg[0] = '\0';
    dinfo.err = jpeg_std_error(&jerr.mgr);
    jerr.mgr.output_message = jpeg_err_output_msg;
    jerr.mgr.error_exit = jpeg_err_exit;
    if (setjmp(jerr.setjmp_buf)) {
        goto bailout;
    }
    jpeg_create_decompress(&dinfo);
    jpeg_mem_src(&dinfo, img, img_size);
    if (jerr.mgr.num_warnings > 0) {
        goto bailout;
    }
    // 读取header后，得到图片color_space和宽高信息，校验一下
    if (jpeg_read_header(&dinfo, TRUE) != JPEG_HEADER_OK) {
        goto bailout;
    }
    // 根据options设置各种dinfo
    if (options != NULL) {
        dinfo.dct_method = options->dct_method;
        dinfo.two_pass_quantize = options->two_pass_quantize;
        dinfo.dither_mode = options->dither_mode;
        dinfo.desired_number_of_colors = options->desired_number_of_colors;
        dinfo.do_fancy_upsampling = options->do_fancy_upsampling;
        if (options->scale_num > 0 && options->scale_denom > 0) {
            dinfo.scale_num = options->scale_num;
            dinfo.scale_denom = options->scale_denom;
        }
    }
    if (dinfo.jpeg_color_space != JCS_GRAYSCALE && dinfo.jpeg_color_space != JCS_YCbCr) {
        snprintf(jerr.last_msg, JMSG_LENGTH_MAX, "unsupported color space, which is %d", dinfo.jpeg_color_space);
        goto bailout;
    }
    // 开始解码图片
    if (jpeg_start_decompress(&dinfo) == FALSE) {
        goto bailout;
    }
    if (options != NULL && options->crop.width > 0 && options->crop.height > 0) {
        // 有图片剪裁的情况，校验输入的crop_width, crop_height是否正确
        if (options->crop.left >= dinfo.image_width || options->crop.top >= dinfo.image_height) {
            goto bailout;
        }
        // 校准width和height，保证不超出图片范围
        if (options->crop.left + options->crop.width > dinfo.image_width) {
            crop_width = dinfo.image_width - options->crop.left;
        } else {
            crop_width = options->crop.width;
        }
        if (options->crop.top + options->crop.height > dinfo.image_height) {
            crop_height = dinfo.image_height - options->crop.top;
        } else {
            crop_height = options->crop.height;
        }
        img_decoded = (JSAMPROW)malloc(sizeof(JSAMPLE) * crop_width * crop_height * dinfo.num_components);
        if (img_decoded == NULL) {
            goto bailout;
        }
        real_left = (JDIMENSION)options->crop.left;
        real_width = (JDIMENSION)crop_width;
        // 需要局部解码图片的话，使用real_left和real_width，因为解码必须整个MCU操作，最终的出来的行还需要一次拷贝才完整。
        if (options->crop.left > 0 || crop_width < dinfo.image_width) {
            jpeg_crop_scanline(&dinfo, &real_left, &real_width);
        }
        // 纵向跳过指定行数
        if (options->crop.top > 0 && (tmp = jpeg_skip_scanlines(&dinfo, (JDIMENSION)options->crop.top)) != options->crop.top) {
            snprintf(jerr.last_msg, JMSG_LENGTH_MAX, "jpeg_skip_scanlines() return %u rather than %u", tmp, options->crop.top);
            goto bailout;
        }
        // 逐行读取scanlines，每行结果用img_row来接，因为MCU只能整个解码，实际real_width有可能比crop_width大。
        img_row = (JSAMPROW)malloc(sizeof(JSAMPLE) * real_width * dinfo.num_components);
        img_row_start = img_row + (sizeof(JSAMPLE) * (options->crop.left - real_left) * dinfo.num_components);
        img_row_size = sizeof(JSAMPLE) * crop_width * dinfo.num_components;
        current_row = img_decoded; // 指向当前第一行的指针
        while (dinfo.output_scanline < options->crop.top + crop_height) {
            // 每次只读一行，因为每行的前面有(options->crop.left-real_left)个像素被剪裁了
            jpeg_read_scanlines(&dinfo, &img_row, 1);
            // 实际上读出来的scanlines会多于需要的像素，所以复制一下到img_decoded
            memcpy(current_row, img_row_start, img_row_size);
            //每读完一行，current_row就移动到下一行
            current_row += img_row_size;
        }
    } else {
        // 无图片剪裁的情况，直接一个buffer copy过去，解码更快。
        img_decoded = (JSAMPROW)malloc(sizeof(JSAMPLE) * dinfo.output_width * dinfo.output_height * dinfo.num_components);
        if (img_decoded == NULL) {
            goto bailout;
        }
        crop_width = dinfo.output_width;
        crop_height = dinfo.output_height;
        // 初始化row_pointer
        row_pointer = (JSAMPARRAY)malloc(sizeof(JSAMPROW) * dinfo.output_height);
        if (row_pointer == NULL) {
            goto bailout;
        }
        img_row_size = sizeof(JSAMPLE) * dinfo.output_width * dinfo.num_components;
        for (tmp = 0; tmp < dinfo.output_height; tmp++) {
            row_pointer[tmp] = &img_decoded[tmp * img_row_size];
        }
        while (dinfo.output_scanline < dinfo.output_height) {
            jpeg_read_scanlines(&dinfo, &row_pointer[dinfo.output_scanline], dinfo.output_height - dinfo.output_scanline);
        }
        jpeg_finish_decompress(&dinfo);
    }
bailout:
    jres->img = img_decoded;
    if (jres->img != NULL) {
        jres->img_size = sizeof(JSAMPLE) * crop_width * crop_height * dinfo.num_components;
    }
    // 如果last_msg非空，从c的栈copy去堆上
    if (jerr.last_msg[0] != '\0') {
        jres->err = malloc(sizeof(char) * JMSG_LENGTH_MAX);
        memcpy(jres->err, jerr.last_msg, JMSG_LENGTH_MAX);
    }
    jres->image_width = crop_width;
    jres->image_height = crop_height;
    jres->origin_width = dinfo.image_width;
    jres->origin_height = dinfo.image_height;
    jres->color_space = dinfo.jpeg_color_space;
    jres->num_components = dinfo.num_components;
    jpeg_destroy_decompress(&dinfo);
    if (img_row != NULL) {
        free(img_row);
    }
    if (row_pointer != NULL) {
        free(row_pointer);
    }
}

// 编码jpeg图片
void jpeg_encode(unsigned char* img, int width, int height, int pixel_format, jpeg_encode_options* options,
    jpeg_encode_result *jres) {
    int            quality    = DEFAULT_QUALITY;
    int            flag       = 0;
    int            sub_sample = TJSAMP_420;
    tjhandle       tj_handler = NULL;

    tj_handler = tjInitCompress();
    if (tj_handler == NULL) {
        goto bailout;
    }
    if (options != NULL) {
        if (options->quality > 0) {
            quality = options->quality;
        }
        flag = options->tj_flag;
        if (options->sub_sample >= 0) {
            sub_sample = options->sub_sample;
        }
    }
    if (tjCompress2(tj_handler, img, width, 0, height, pixel_format, &(jres->img), &(jres->img_size), sub_sample,
        quality, flag) < 0) {
        goto bailout;
    }
    tjDestroy(tj_handler);
    return;
bailout:
    // 错误异常处理
    jres->err = (char*)malloc(sizeof(char) * JMSG_LENGTH_MAX);
    memcpy(jres->err, tjGetErrorStr2(tj_handler), JMSG_LENGTH_MAX);
    if (tj_handler != NULL) {
        tjDestroy(tj_handler);
    }
}
