use actix_web;


#[actix_web::get("/test_plain")]
async fn test_handler() -> impl actix_web::Responder {
    actix_web::HttpResponse::Ok().body("Hello world!")
}

//It would be great idea to add SIGTERM handler and exit gracefully
//https://stackoverflow.com/questions/26280859/how-to-catch-signals-in-rust
//But for now i'll just `kill -9`
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    actix_web::HttpServer::new(|| {
        actix_web::App::new()
            .service(test_handler)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
