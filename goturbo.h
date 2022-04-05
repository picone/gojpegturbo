#ifndef GO_JPEG_TURBO
#define GO_JPEG_TURBO

#include <stdio.h>
#include <memory.h>
#include <stdlib.h>
#include <setjmp.h>
#include "turbojpeg.h"
#include "jpeglib.h"

#define DEFAULT_QUALITY 95

// 搞一个新的err mgr，因为原来的不能保存last_msg信息。
typedef struct my_jpeg_err_mgr {
    struct jpeg_error_mgr mgr;
    jmp_buf setjmp_buf;
    char last_msg[JMSG_LENGTH_MAX];
} my_jpeg_err_mgr;

typedef struct crop_rect {
    unsigned int left;
    unsigned int top;
    unsigned int width;
    unsigned int height;
} crop_rect;

typedef struct jpeg_decode_options {
    struct crop_rect crop;
    J_DCT_METHOD dct_method;
    boolean two_pass_quantize;
    J_DITHER_MODE dither_mode;
    int desired_number_of_colors;
    boolean do_fancy_upsampling;
    unsigned int scale_num, scale_denom;
} jpeg_decode_options;

typedef struct jpeg_decode_result {
    unsigned char* img;
    unsigned int img_size;
    unsigned int image_width;
    unsigned int image_height;
    unsigned int origin_width;
    unsigned int origin_height;
    J_COLOR_SPACE color_space;
    int num_components;
    char* err;
} jpeg_decode_result;

typedef struct jpeg_encode_options {
    int quality;
    int tj_flag;
    int sub_sample;
} jpeg_encode_options;

typedef struct jpeg_encode_result {
    unsigned char* img;
    unsigned long img_size;
    char* err;
} jpeg_encode_result;

// 覆盖原来的output_message方法，因为原来的会打印到控制台。
static void jpeg_err_output_msg(j_common_ptr cinfo);

// 覆盖原来的error_exit方法，因为原来的错误会调用exit函数导致进程退出。
static void jpeg_err_exit(j_common_ptr cinfo);

// 解码jpeg图片
void jpeg_decode(unsigned char* img, unsigned int img_size, jpeg_decode_options* options, jpeg_decode_result* jres);

// 编码jpeg图片
void jpeg_encode(unsigned char* img, int width, int height, int pixel_format, jpeg_encode_options* options,
    jpeg_encode_result *jres);

#endif