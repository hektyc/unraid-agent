#!/usr/bin/env python3
"""Generate the unraid-agent AI robot icon.

Design: modern flat robot head on a purple->blue diagonal gradient
rounded-square background, orange robot with glowing eyes.
Rendered at 1024px and downsampled for crisp anti-aliased edges.
"""

from PIL import Image, ImageDraw, ImageFilter
import math
import os

S = 1024  # render canvas
OUT = os.path.join(os.path.dirname(__file__), "..", "plugin", "images")


def lerp(a, b, t):
    return tuple(int(a[i] + (b[i] - a[i]) * t) for i in range(len(a)))


def diagonal_gradient(size, c1, c2, radius):
    """Rounded-square with a diagonal linear gradient."""
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    grad = Image.new("RGBA", (size, size))
    px = grad.load()
    for y in range(size):
        for x in range(size):
            t = (x + y) / (2 * size)
            px[x, y] = (*lerp(c1, c2, t), 255)
    mask = Image.new("L", (size, size), 0)
    d = ImageDraw.Draw(mask)
    d.rounded_rectangle([0, 0, size - 1, size - 1], radius=radius, fill=255)
    img.paste(grad, (0, 0), mask)
    return img, mask


def glow(canvas, center, r, color, alpha=90):
    layer = Image.new("RGBA", canvas.size, (0, 0, 0, 0))
    d = ImageDraw.Draw(layer)
    d.ellipse([center[0] - r, center[1] - r, center[0] + r, center[1] + r],
              fill=(*color, alpha))
    layer = layer.filter(ImageFilter.GaussianBlur(r * 0.6))
    canvas.alpha_composite(layer)


def main():
    # palette
    PURPLE = (58, 24, 110)     # deep violet
    BLUE = (30, 90, 200)       # modern blue
    ORANGE = (249, 115, 22)    # modern orange
    ORANGE_HI = (253, 154, 60) # lighter orange
    FACE = (255, 237, 213)     # warm off-white faceplate
    EYE = (67, 56, 202)        # indigo eyes
    EYE_GLOW = (129, 140, 248) # periwinkle glow

    img, _ = diagonal_gradient(S, PURPLE, BLUE, int(S * 0.22))

    # subtle inner border highlight
    d = ImageDraw.Draw(img)
    d.rounded_rectangle([S*0.015, S*0.015, S*0.985, S*0.985],
                        radius=int(S*0.21), outline=(255, 255, 255, 38),
                        width=int(S*0.008))

    cx = S / 2

    # ---------- antenna ----------
    ant_w = int(S * 0.022)
    ant_top = int(S * 0.085)
    ant_bot = int(S * 0.245)
    d.rounded_rectangle([cx - ant_w/2, ant_top, cx + ant_w/2, ant_bot],
                        radius=ant_w//2, fill=ORANGE_HI)
    tip_r = int(S * 0.045)
    glow(img, (cx, ant_top), int(tip_r * 1.9), ORANGE_HI, 70)
    d = ImageDraw.Draw(img)
    d.ellipse([cx - tip_r, ant_top - tip_r, cx + tip_r, ant_top + tip_r],
              fill=ORANGE_HI)

    # ---------- ears ----------
    ear_w = int(S * 0.055)
    ear_h = int(S * 0.135)
    ear_y = int(S * 0.415)
    head_hw = int(S * 0.31)  # head half width
    for side in (-1, 1):
        x0 = cx + side * head_hw - (ear_w if side < 0 else 0)
        x1 = x0 + ear_w
        d.rounded_rectangle([x0, ear_y, x1, ear_y + ear_h],
                            radius=int(S*0.03), fill=ORANGE)

    # ---------- head ----------
    hx0, hy0 = cx - head_hw, int(S * 0.245)
    hx1, hy1 = cx + head_hw, int(S * 0.79)
    head = Image.new("RGBA", (S, S), (0, 0, 0, 0))
    hd = ImageDraw.Draw(head)
    # vertical orange gradient on the head
    hg = Image.new("RGBA", (S, S))
    hp = hg.load()
    for y in range(S):
        t = min(1, max(0, (y - hy0) / (hy1 - hy0)))
        for x in range(S):
            hp[x, y] = (*lerp(ORANGE_HI, ORANGE, t), 255)
    hmask = Image.new("L", (S, S), 0)
    hd = ImageDraw.Draw(hmask)
    hd.rounded_rectangle([hx0, hy0, hx1, hy1], radius=int(S*0.14), fill=255)
    head.paste(hg, (0, 0), hmask)
    img.alpha_composite(head)

    # ---------- faceplate ----------
    fx0, fy0 = cx - int(S*0.24), int(S * 0.335)
    fx1, fy1 = cx + int(S*0.24), int(S * 0.70)
    d = ImageDraw.Draw(img)
    d.rounded_rectangle([fx0, fy0, fx1, fy1], radius=int(S*0.10), fill=FACE)

    # ---------- eyes ----------
    eye_w = int(S * 0.075)
    eye_h = int(S * 0.115)
    eye_y = int(S * 0.43)
    eye_dx = int(S * 0.115)
    for side in (-1, 1):
        ex = cx + side * eye_dx
        glow(img, (ex, eye_y + eye_h//2), int(eye_w * 1.5), EYE_GLOW, 85)
    d = ImageDraw.Draw(img)
    for side in (-1, 1):
        ex = cx + side * eye_dx
        d.rounded_rectangle([ex - eye_w/2, eye_y, ex + eye_w/2, eye_y + eye_h],
                            radius=int(eye_w * 0.45), fill=EYE)
        # eye highlight
        hl = int(eye_w * 0.22)
        d.ellipse([ex - eye_w*0.18, eye_y + eye_h*0.16,
                   ex - eye_w*0.18 + hl, eye_y + eye_h*0.16 + hl],
                  fill=(255, 255, 255, 220))

    # ---------- mouth (speaker grille) ----------
    grille_w = int(S * 0.19)
    grille_y = int(S * 0.615)
    bar_w = int(S * 0.018)
    bar_h = int(S * 0.045)
    gap = int(S * 0.028)
    n = 3
    total = n * bar_w + (n - 1) * gap
    gx = cx - total / 2
    for i in range(n):
        d.rounded_rectangle([gx + i * (bar_w + gap), grille_y,
                             gx + i * (bar_w + gap) + bar_w, grille_y + bar_h],
                            radius=bar_w // 2, fill=(234, 88, 12))

    # ---------- chin shadow accent ----------
    d.rounded_rectangle([cx - int(S*0.16), int(S*0.735), cx + int(S*0.16), int(S*0.76)],
                        radius=int(S*0.012), fill=(234, 88, 12, 140))

    os.makedirs(OUT, exist_ok=True)
    for size, name in [(128, "unraid-agent-128.png"),
                       (48, "unraid-agent-48.png"),
                       (128, "unraid-agent.png")]:
        out = img.resize((size, size), Image.LANCZOS)
        out.save(os.path.join(OUT, name), optimize=True)
        print(f"wrote {name} ({size}x{size})")

    # 16px settings-menu icon: icons/<lowercase-title-no-spaces>.png
    icons_dir = os.path.join(OUT, "..", "icons")
    os.makedirs(icons_dir, exist_ok=True)
    img.resize((16, 16), Image.LANCZOS).save(
        os.path.join(icons_dir, "unraidagent.png"), optimize=True)
    print("wrote icons/unraidagent.png (16x16)")

    # preview at 256 for inspection
    img.resize((256, 256), Image.LANCZOS).save("/tmp/unraid-agent-preview.png")


if __name__ == "__main__":
    main()
