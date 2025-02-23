//extern crate rocket;
use rocket;

#[rocket::get("/test_plain")]
fn test_handler() -> &'static str {
    return "Hello world!";
}

#[rocket::launch]
fn rocket_start() -> _ {
    println!("Hello, world!");
    rocket::build().mount("/", rocket::routes![test_handler])
}

