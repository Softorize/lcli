#!/usr/bin/env python3
"""
LinkedIn Carousel Generator - Creates slide images and combines into PDF.

Usage:
    python3 carousel_generator.py slides.json [output.pdf]

slides.json format:
{
    "brand": {
        "name": "Elixion.ai",
        "accent_color": "#0077B5",
        "bg_color": "#0F0F0F",
        "text_color": "#FFFFFF",
        "subtitle_color": "#B0B0B0",
        "author": "Zhirayr Gumruyan",
        "handle": "@zhirayrgumruyan"
    },
    "slides": [
        {
            "type": "cover",
            "title": "Big Bold Title",
            "subtitle": "Supporting text"
        },
        {
            "type": "content",
            "title": "Slide Title",
            "body": "Main content text",
            "highlight": "Optional highlighted phrase"
        },
        {
            "type": "list",
            "title": "Key Points",
            "items": ["Point 1", "Point 2", "Point 3"]
        },
        {
            "type": "quote",
            "quote": "The quote text",
            "attribution": "- Author"
        },
        {
            "type": "cta",
            "title": "Call to Action",
            "body": "What do you think?",
            "cta": "Follow for more insights"
        }
    ]
}
"""

import json
import sys
import textwrap
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont
from reportlab.lib.pagesizes import letter
from reportlab.pdfgen import canvas as pdf_canvas


# --- Constants ---
SLIDE_W, SLIDE_H = 1080, 1350  # LinkedIn optimal carousel size (4:5 ratio)
MARGIN = 80
CONTENT_W = SLIDE_W - 2 * MARGIN


# --- Font helpers ---
def load_font(bold=False, size=40):
    """Load system fonts with fallback."""
    candidates = []
    if bold:
        candidates = [
            "/System/Library/Fonts/Supplemental/Arial Bold.ttf",
            "/System/Library/Fonts/Helvetica.ttc",
            "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
        ]
    else:
        candidates = [
            "/System/Library/Fonts/Supplemental/Arial.ttf",
            "/System/Library/Fonts/Helvetica.ttc",
            "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
        ]
    for path in candidates:
        try:
            return ImageFont.truetype(path, size)
        except (OSError, IOError):
            continue
    return ImageFont.load_default()


FONT_TITLE = load_font(bold=True, size=56)
FONT_TITLE_LARGE = load_font(bold=True, size=64)
FONT_BODY = load_font(bold=False, size=36)
FONT_BODY_BOLD = load_font(bold=True, size=36)
FONT_SUBTITLE = load_font(bold=False, size=32)
FONT_SMALL = load_font(bold=False, size=24)
FONT_SMALL_BOLD = load_font(bold=True, size=24)
FONT_QUOTE = load_font(bold=True, size=44)
FONT_NUMBER = load_font(bold=True, size=80)
FONT_CTA_TITLE = load_font(bold=True, size=52)


def hex_to_rgb(hex_color):
    """Convert hex color to RGB tuple."""
    h = hex_color.lstrip("#")
    return tuple(int(h[i : i + 2], 16) for i in (0, 2, 4))


def draw_rounded_rect(draw, xy, radius, fill):
    """Draw a rounded rectangle."""
    x0, y0, x1, y1 = xy
    draw.rectangle([x0 + radius, y0, x1 - radius, y1], fill=fill)
    draw.rectangle([x0, y0 + radius, x1, y1 - radius], fill=fill)
    draw.pieslice([x0, y0, x0 + 2 * radius, y0 + 2 * radius], 180, 270, fill=fill)
    draw.pieslice([x1 - 2 * radius, y0, x1, y0 + 2 * radius], 270, 360, fill=fill)
    draw.pieslice([x0, y1 - 2 * radius, x0 + 2 * radius, y1], 90, 180, fill=fill)
    draw.pieslice([x1 - 2 * radius, y1 - 2 * radius, x1, y1], 0, 90, fill=fill)


def wrap_text(text, font, max_width, draw):
    """Wrap text to fit within max_width pixels."""
    words = text.split()
    lines = []
    current_line = ""

    for word in words:
        test_line = f"{current_line} {word}".strip()
        bbox = draw.textbbox((0, 0), test_line, font=font)
        if bbox[2] - bbox[0] <= max_width:
            current_line = test_line
        else:
            if current_line:
                lines.append(current_line)
            current_line = word

    if current_line:
        lines.append(current_line)

    return lines


def draw_slide_number(draw, slide_num, total, brand):
    """Draw slide indicator dots at the bottom."""
    dot_radius = 6
    dot_spacing = 24
    total_width = total * dot_spacing
    start_x = (SLIDE_W - total_width) // 2
    y = SLIDE_H - 50

    accent = hex_to_rgb(brand["accent_color"])
    dim = hex_to_rgb(brand.get("subtitle_color", "#666666"))

    for i in range(total):
        x = start_x + i * dot_spacing + dot_radius
        color = accent if i == slide_num else dim
        draw.ellipse([x - dot_radius, y - dot_radius, x + dot_radius, y + dot_radius], fill=color)


def draw_branding(draw, brand, light=False):
    """Draw brand name at top-right corner."""
    color = hex_to_rgb(brand.get("subtitle_color", "#B0B0B0"))
    text = brand.get("name", "")
    if text:
        bbox = draw.textbbox((0, 0), text, font=FONT_SMALL_BOLD)
        tw = bbox[2] - bbox[0]
        draw.text((SLIDE_W - MARGIN - tw, 40), text, fill=color, font=FONT_SMALL_BOLD)


def draw_gradient_bar(draw, y, width, accent_color):
    """Draw a horizontal accent bar."""
    accent = hex_to_rgb(accent_color)
    draw.rectangle([MARGIN, y, MARGIN + width, y + 4], fill=accent)


# --- Slide renderers ---

def render_cover(img, draw, slide, brand, slide_num, total):
    """Render cover slide with big title and subtitle."""
    bg = hex_to_rgb(brand["bg_color"])
    accent = hex_to_rgb(brand["accent_color"])
    text_color = hex_to_rgb(brand["text_color"])
    sub_color = hex_to_rgb(brand.get("subtitle_color", "#B0B0B0"))

    # Accent bar at top
    draw.rectangle([0, 0, SLIDE_W, 8], fill=accent)

    # Brand
    draw_branding(draw, brand)

    # Title - centered vertically
    title = slide.get("title", "")
    lines = wrap_text(title, FONT_TITLE_LARGE, CONTENT_W, draw)
    line_height = 78
    total_text_height = len(lines) * line_height
    start_y = (SLIDE_H - total_text_height) // 2 - 60

    for i, line in enumerate(lines):
        draw.text((MARGIN, start_y + i * line_height), line, fill=text_color, font=FONT_TITLE_LARGE)

    # Accent underline below title
    underline_y = start_y + len(lines) * line_height + 20
    draw.rectangle([MARGIN, underline_y, MARGIN + 120, underline_y + 5], fill=accent)

    # Subtitle
    subtitle = slide.get("subtitle", "")
    if subtitle:
        sub_lines = wrap_text(subtitle, FONT_SUBTITLE, CONTENT_W, draw)
        sub_start = underline_y + 40
        for i, line in enumerate(sub_lines):
            draw.text((MARGIN, sub_start + i * 44), line, fill=sub_color, font=FONT_SUBTITLE)

    # Author at bottom
    author = brand.get("author", "")
    handle = brand.get("handle", "")
    if author:
        draw.text((MARGIN, SLIDE_H - 130), author, fill=text_color, font=FONT_BODY_BOLD)
    if handle:
        draw.text((MARGIN, SLIDE_H - 90), handle, fill=sub_color, font=FONT_SMALL)

    draw_slide_number(draw, slide_num, total, brand)


def render_content(img, draw, slide, brand, slide_num, total):
    """Render a content slide with title and body text."""
    accent = hex_to_rgb(brand["accent_color"])
    text_color = hex_to_rgb(brand["text_color"])
    sub_color = hex_to_rgb(brand.get("subtitle_color", "#B0B0B0"))

    draw_branding(draw, brand)

    # Title
    title = slide.get("title", "")
    title_lines = wrap_text(title, FONT_TITLE, CONTENT_W, draw)
    y = 120
    for line in title_lines:
        draw.text((MARGIN, y), line, fill=text_color, font=FONT_TITLE)
        y += 68

    # Accent bar
    draw_gradient_bar(draw, y + 10, 80, brand["accent_color"])
    y += 50

    # Body text
    body = slide.get("body", "")
    body_lines = wrap_text(body, FONT_BODY, CONTENT_W, draw)
    for line in body_lines:
        draw.text((MARGIN, y), line, fill=sub_color, font=FONT_BODY)
        y += 50

    # Highlight box (if present)
    highlight = slide.get("highlight", "")
    if highlight:
        y += 30
        hl_lines = wrap_text(highlight, FONT_BODY_BOLD, CONTENT_W - 60, draw)
        box_height = len(hl_lines) * 50 + 40
        draw_rounded_rect(draw, (MARGIN, y, SLIDE_W - MARGIN, y + box_height), 16, fill=hex_to_rgb("#1A1A2E"))
        draw.rectangle([MARGIN, y, MARGIN + 5, y + box_height], fill=accent)  # left accent bar
        hy = y + 20
        for line in hl_lines:
            draw.text((MARGIN + 30, hy), line, fill=hex_to_rgb("#FFFFFF"), font=FONT_BODY_BOLD)
            hy += 50

    draw_slide_number(draw, slide_num, total, brand)


def render_list(img, draw, slide, brand, slide_num, total):
    """Render a list slide with numbered items."""
    accent = hex_to_rgb(brand["accent_color"])
    text_color = hex_to_rgb(brand["text_color"])
    sub_color = hex_to_rgb(brand.get("subtitle_color", "#B0B0B0"))

    draw_branding(draw, brand)

    # Title
    title = slide.get("title", "")
    title_lines = wrap_text(title, FONT_TITLE, CONTENT_W, draw)
    y = 120
    for line in title_lines:
        draw.text((MARGIN, y), line, fill=text_color, font=FONT_TITLE)
        y += 68

    draw_gradient_bar(draw, y + 10, 80, brand["accent_color"])
    y += 60

    # List items
    items = slide.get("items", [])
    for idx, item in enumerate(items, 1):
        # Number circle
        circle_x = MARGIN + 28
        circle_y = y + 22
        draw.ellipse(
            [circle_x - 28, circle_y - 28, circle_x + 28, circle_y + 28],
            fill=accent,
        )
        num_text = str(idx)
        num_bbox = draw.textbbox((0, 0), num_text, font=FONT_BODY_BOLD)
        num_w = num_bbox[2] - num_bbox[0]
        num_h = num_bbox[3] - num_bbox[1]
        draw.text(
            (circle_x - num_w // 2, circle_y - num_h // 2 - 4),
            num_text,
            fill=hex_to_rgb("#FFFFFF"),
            font=FONT_BODY_BOLD,
        )

        # Item text (wrapped)
        item_lines = wrap_text(item, FONT_BODY, CONTENT_W - 90, draw)
        text_x = MARGIN + 72
        for il, iline in enumerate(item_lines):
            draw.text((text_x, y + il * 46), iline, fill=sub_color, font=FONT_BODY)

        y += max(len(item_lines) * 46, 56) + 30

    draw_slide_number(draw, slide_num, total, brand)


def render_quote(img, draw, slide, brand, slide_num, total):
    """Render a quote slide with large quotation marks."""
    accent = hex_to_rgb(brand["accent_color"])
    text_color = hex_to_rgb(brand["text_color"])
    sub_color = hex_to_rgb(brand.get("subtitle_color", "#B0B0B0"))

    draw_branding(draw, brand)

    # Large quotation mark
    quote_mark_font = load_font(bold=True, size=200)
    draw.text((MARGIN - 10, 150), "\u201C", fill=accent, font=quote_mark_font)

    # Quote text
    quote = slide.get("quote", "")
    y = 380
    quote_lines = wrap_text(quote, FONT_QUOTE, CONTENT_W, draw)
    for line in quote_lines:
        draw.text((MARGIN, y), line, fill=text_color, font=FONT_QUOTE)
        y += 60

    # Attribution
    attr = slide.get("attribution", "")
    if attr:
        y += 30
        draw.text((MARGIN, y), attr, fill=sub_color, font=FONT_SUBTITLE)

    draw_slide_number(draw, slide_num, total, brand)


def render_cta(img, draw, slide, brand, slide_num, total):
    """Render a call-to-action closing slide."""
    accent = hex_to_rgb(brand["accent_color"])
    text_color = hex_to_rgb(brand["text_color"])
    sub_color = hex_to_rgb(brand.get("subtitle_color", "#B0B0B0"))

    # Accent bar at top
    draw.rectangle([0, 0, SLIDE_W, 8], fill=accent)

    draw_branding(draw, brand)

    # Title
    title = slide.get("title", "")
    title_lines = wrap_text(title, FONT_CTA_TITLE, CONTENT_W, draw)
    y = (SLIDE_H // 2) - 150
    for line in title_lines:
        draw.text((MARGIN, y), line, fill=text_color, font=FONT_CTA_TITLE)
        y += 66

    # Body
    body = slide.get("body", "")
    if body:
        y += 20
        body_lines = wrap_text(body, FONT_BODY, CONTENT_W, draw)
        for line in body_lines:
            draw.text((MARGIN, y), line, fill=sub_color, font=FONT_BODY)
            y += 50

    # CTA button
    cta = slide.get("cta", "")
    if cta:
        y += 40
        cta_bbox = draw.textbbox((0, 0), cta, font=FONT_BODY_BOLD)
        cta_w = cta_bbox[2] - cta_bbox[0] + 60
        cta_h = 64
        cta_x = MARGIN
        draw_rounded_rect(draw, (cta_x, y, cta_x + cta_w, y + cta_h), 12, fill=accent)
        draw.text((cta_x + 30, y + 12), cta, fill=hex_to_rgb("#FFFFFF"), font=FONT_BODY_BOLD)

    # Author at bottom
    author = brand.get("author", "")
    handle = brand.get("handle", "")
    if author:
        draw.text((MARGIN, SLIDE_H - 130), author, fill=text_color, font=FONT_BODY_BOLD)
    if handle:
        draw.text((MARGIN, SLIDE_H - 90), handle, fill=sub_color, font=FONT_SMALL)

    draw_slide_number(draw, slide_num, total, brand)


RENDERERS = {
    "cover": render_cover,
    "content": render_content,
    "list": render_list,
    "quote": render_quote,
    "cta": render_cta,
}


def generate_slide(slide_data, brand, slide_num, total_slides):
    """Generate a single slide image."""
    bg = hex_to_rgb(brand["bg_color"])
    img = Image.new("RGB", (SLIDE_W, SLIDE_H), color=bg)
    draw = ImageDraw.Draw(img)

    slide_type = slide_data.get("type", "content")
    renderer = RENDERERS.get(slide_type, render_content)
    renderer(img, draw, slide_data, brand, slide_num, total_slides)

    return img


def images_to_pdf(images, output_path):
    """Combine slide images into a PDF for LinkedIn carousel upload."""
    c = pdf_canvas.Canvas(str(output_path))

    for img in images:
        # Save temp image
        tmp_path = Path(output_path).parent / "_tmp_slide.png"
        img.save(str(tmp_path), "PNG")

        # Set page size to match image aspect ratio (scale to reasonable PDF size)
        pdf_w = 540  # points (7.5 inches)
        pdf_h = int(pdf_w * (SLIDE_H / SLIDE_W))
        c.setPageSize((pdf_w, pdf_h))
        c.drawImage(str(tmp_path), 0, 0, width=pdf_w, height=pdf_h)
        c.showPage()

        tmp_path.unlink(missing_ok=True)

    c.save()


def main():
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} slides.json [output.pdf]")
        sys.exit(1)

    input_path = Path(sys.argv[1])
    output_path = Path(sys.argv[2]) if len(sys.argv) > 2 else input_path.with_suffix(".pdf")

    with open(input_path) as f:
        data = json.load(f)

    brand = data.get("brand", {})
    slides = data.get("slides", [])
    total = len(slides)

    print(f"Generating {total} slides...")

    images = []
    slide_dir = output_path.parent / "slides"
    slide_dir.mkdir(exist_ok=True)

    for i, slide in enumerate(slides):
        img = generate_slide(slide, brand, i, total)
        images.append(img)

        # Also save individual PNGs
        png_path = slide_dir / f"slide_{i + 1:02d}.png"
        img.save(str(png_path), "PNG")
        print(f"  Saved {png_path}")

    # Generate PDF
    images_to_pdf(images, output_path)
    print(f"\nPDF carousel: {output_path}")
    print(f"Individual slides: {slide_dir}/")
    print("Done!")


if __name__ == "__main__":
    main()
