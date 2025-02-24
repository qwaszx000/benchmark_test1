#![feature(core_intrinsics)]
#![feature(lang_items)]
#![allow(clippy::missing_safety_doc, unused_imports, dead_code)]

use core::intrinsics::likely;
use core::ptr::copy_nonoverlapping;

use faf::epoll;
use faf::util::memcmp;

const GET_METHOD: &[u8] = b"GET";
const GET_LENGTH: usize = GET_METHOD.len();

const ROUTE_PATH: &[u8] = b"/test_plain";
const ROUTE_LENGTH: usize = ROUTE_PATH.len();

const RESPONSE_404: &[u8] = b"HTTP/1.1 404 Not Found\r\nContent-type: text/plain\r\nContent-Length: 9\r\n\r\nNot Found";
const RESPONSE_404_LENGTH: usize = RESPONSE_404.len();

const RESPONSE_200: &[u8] = b"HTTP/1.1 200 Ok\r\nContent-type: text/plain\r\nContent-Length: 12\r\n\r\nHello world!";
const RESPONSE_200_LENGTH: usize = RESPONSE_200.len();

#[inline(always)]
fn test_handler(
    method: *const u8,
    method_len: usize,
    path: *const u8,
    path_len: usize,
    response_buf: *mut u8,
    date_buf: *const u8
) -> usize {
    unsafe {
        //Test by length
        if likely(method_len == GET_LENGTH && path_len == ROUTE_LENGTH) {
            //Test method
            if likely(memcmp(GET_METHOD.as_ptr(), method, GET_LENGTH) == 0) {
                //Test request path
                if likely(memcmp(ROUTE_PATH.as_ptr(), path, ROUTE_LENGTH) == 0) {
                    copy_nonoverlapping(RESPONSE_200.as_ptr(), response_buf, RESPONSE_200_LENGTH);
                    return RESPONSE_200_LENGTH;
                }
            }
        }

        copy_nonoverlapping(RESPONSE_404.as_ptr(), response_buf, RESPONSE_404_LENGTH);
        return RESPONSE_404_LENGTH;
    }
    return 0;
}

#[inline(always)]
fn main() {
    epoll::go(8080, test_handler)
}


