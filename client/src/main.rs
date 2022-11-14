use std::net::TcpStream;
use std::str;
use std::io::{self, BufRead, BufReader, Write};
use std::env;

fn main() {
    let args: Vec<String> = env::args().collect();
    let port = args[1].clone(); 
    connect(port);
   
}

fn connect(port: String) {
    let address = format!("127.0.0.1:{}", port);
    let mut stream = TcpStream::connect(address).expect("Unable to connect proxy");
    loop {
        let mut input = String::new();
        let mut buffer: Vec<u8> = Vec::new();

        io::stdin().read_line(&mut input).expect("error reading from stdin");
        stream.write(input.as_bytes()).expect("error writing to proxy");

        let mut reader = BufReader::new(&stream);
        reader.read_until(b'\n', &mut buffer).expect("error reading into buffer");

        
        print!("{}", str::from_utf8(&buffer).expect("error writing into buffer as string"));
    }
}