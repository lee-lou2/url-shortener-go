use axum::{
    extract::Query,
    http::{header, HeaderMap, StatusCode},
    response::{IntoResponse, Response},
};
use image::{imageops, imageops::FilterType, ImageBuffer, ImageFormat, Rgba};
use qrcode::Color as QrColor;
use qrcode::QrCode;
use serde::Deserialize;
use std::io::Cursor;
use urlencoding::decode;

/// QR code generation request query parameters
#[derive(Deserialize, Debug)]
pub struct GenerateQrQuery {
    url: String,
    size: Option<String>,  // Size (e.g., "512x512")
    color: Option<String>, // Foreground (module) color (e.g., "black", "#FF0000")
    bg: Option<String>,    // Background color (e.g., "white", "#00FF00")
}

/// Parse string in "WxH" format (returns width and height)
fn parse_size(size_str: &str) -> Result<(u32, u32), &'static str> {
    let parts: Vec<&str> = size_str.split('x').collect();
    if parts.len() != 2 {
        return Err("Invalid size format. Use WxH (e.g., 512x512)");
    }
    let width = parts[0].parse::<u32>().map_err(|_| "Invalid width")?;
    let height = parts[1].parse::<u32>().map_err(|_| "Invalid height")?;
    if width == 0 || height == 0 {
        return Err("Width and height must be greater than 0");
    }
    Ok((width, height))
}

/// Parse color string (same as before)
fn parse_color(color_str: &str) -> Result<Rgba<u8>, &'static str> {
    match color_str.to_lowercase().as_str() {
        "black" => Ok(Rgba([0, 0, 0, 255])),
        "white" => Ok(Rgba([255, 255, 255, 255])),
        "red" => Ok(Rgba([255, 0, 0, 255])),
        "green" => Ok(Rgba([0, 128, 0, 255])),
        "blue" => Ok(Rgba([0, 0, 255, 255])),
        "yellow" => Ok(Rgba([255, 255, 0, 255])),
        "pink" => Ok(Rgba([255, 192, 203, 255])),
        "purple" => Ok(Rgba([128, 0, 128, 255])),
        "orange" => Ok(Rgba([255, 165, 0, 255])),
        "brown" => Ok(Rgba([139, 69, 19, 255])),
        "gray" => Ok(Rgba([128, 128, 128, 255])),
        "gold" => Ok(Rgba([255, 215, 0, 255])),
        "silver" => Ok(Rgba([192, 192, 192, 255])),
        "cyan" => Ok(Rgba([0, 255, 255, 255])),
        "lime" => Ok(Rgba([0, 255, 0, 255])),
        "teal" => Ok(Rgba([0, 128, 128, 255])),
        "navy" => Ok(Rgba([0, 0, 128, 255])),
        "maroon" => Ok(Rgba([128, 0, 0, 255])),
        "olive" => Ok(Rgba([128, 128, 0, 255])),
        _ => {
            if color_str.starts_with('#') && color_str.len() == 7 {
                let r = u8::from_str_radix(&color_str[1..3], 16)
                    .map_err(|_| "Invalid hex color red component")?;
                let g = u8::from_str_radix(&color_str[3..5], 16)
                    .map_err(|_| "Invalid hex color green component")?;
                let b = u8::from_str_radix(&color_str[5..7], 16)
                    .map_err(|_| "Invalid hex color blue component")?;
                Ok(Rgba([r, g, b, 255]))
            } else {
                Err("Invalid color format. Use name (e.g., black) or hex (#RRGGBB)")
            }
        }
    }
}

/// Generate QR Code Image Handler (using qrcode-rust)
/// GET /v1/qr/?url=...&size=...&color=...&bg=...
/// Response: QR code image in image/png format
pub async fn generate_qr_handler(
    Query(query): Query<GenerateQrQuery>,
) -> Result<Response, (StatusCode, String)> {
    // 1. Decode the URL parameter
    let decoded_url = decode(&query.url)
        .map_err(|e| {
            (
                StatusCode::BAD_REQUEST,
                format!("Invalid URL encoding: {}", e),
            )
        })?
        .into_owned();

    // 2. Process size parameter
    let (final_width, final_height) = query
        .size
        .as_deref()
        .map(parse_size)
        .transpose()
        .map_err(|e| (StatusCode::BAD_REQUEST, e.to_string()))?
        .unwrap_or((256, 256));

    // 3. Process color (foreground) parameter
    let foreground_color = query
        .color
        .as_deref()
        .map(parse_color)
        .transpose()
        .map_err(|e| (StatusCode::BAD_REQUEST, e.to_string()))?
        .unwrap_or(Rgba([0, 0, 0, 255]));

    // 4. Process bg (background) parameter
    let background_color = query
        .bg
        .as_deref()
        .map(parse_color)
        .transpose()
        .map_err(|e| (StatusCode::BAD_REQUEST, e.to_string()))?
        .unwrap_or(Rgba([255, 255, 255, 0]));

    // 5. Generate QR code data
    let code = QrCode::new(&decoded_url).map_err(|e| {
        (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("Failed to create QR code data: {}", e),
        )
    })?;

    // 6. Create initial QR code image buffer (without padding)
    let qrcode_module_width = code.width() as u32;
    let module_colors = code.to_colors();
    let mut qr_image_buffer_rgba = ImageBuffer::new(qrcode_module_width, qrcode_module_width);

    for y in 0..qrcode_module_width {
        for x in 0..qrcode_module_width {
            let index = (y * qrcode_module_width + x) as usize;
            let pixel_color = match module_colors[index] {
                QrColor::Dark => foreground_color,
                QrColor::Light => background_color,
            };
            qr_image_buffer_rgba.put_pixel(x, y, pixel_color);
        }
    }

    // 7. Create a larger buffer with padding and draw the QR code onto it
    let padding = 1;
    let padded_size = qrcode_module_width + padding * 2;
    let mut padded_qr_buffer = ImageBuffer::from_pixel(padded_size, padded_size, background_color);

    imageops::overlay(
        &mut padded_qr_buffer,
        &qr_image_buffer_rgba,
        padding.into(),
        padding.into(),
    );

    // 8. Resize the padded image to the final requested size
    let final_image = imageops::resize(
        &padded_qr_buffer,
        final_width,
        final_height,
        FilterType::Nearest,
    );

    // 9. Encode image data in PNG format
    let mut buf = Cursor::new(Vec::new());
    final_image
        .write_to(&mut buf, ImageFormat::Png)
        .map_err(|e| {
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                format!("Failed to encode PNG: {}", e),
            )
        })?;

    // 10. Create HTTP response
    let mut headers = HeaderMap::new();
    headers.insert(header::CONTENT_TYPE, "image/png".parse().unwrap());

    Ok((headers, buf.into_inner()).into_response())
}
