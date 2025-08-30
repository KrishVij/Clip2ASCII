Clip2ASCII 
![imageIcon](https://github.com/user-attachments/assets/6839b57b-2109-4ffc-acda-62b8593222b8)

Clip2ASCII is a Go-based desktop application that leverages FFmpeg for video processing and custom pixel-to-ASCII mapping to create retro-style ASCII animations.

📸 Demo

https://github.com/user-attachments/assets/4f3e7cc4-3604-4f5f-9588-fcd916262972

Features

Video to ASCII Conversion: Transform short videos (up to 30 seconds) into vibrant, colored ASCII art videos.
Intuitive GUI: Built with Ebiten for a smooth, cross-platform desktop experience. Select videos, preview thumbnails, and track progress with a sleek UI.
High-Quality Processing:

Extracts frames using FFmpeg.
Resizes and gamma-corrects pixels for accurate brightness mapping.
Maps brightness to a custom ASCII charset (from darkest . to brightest @).
Preserves colors in the final ASCII frames.


Automatic Cleanup: Deletes temporary frame folders after conversion.
Output: Saves the ASCII video in your home directory with the original filename.

🛠️ Installation
Prerequisites

Go: Version 1.20+ (install from golang.org).
FFmpeg: Required for video frame extraction and stitching. Install via:

Windows: Download from ffmpeg.org and add to PATH.
Link: [https://ffmpeg.org/](#FFmpeg)

Font File: Place a Font.ttf (e.g., a monospace font like Consolas or Courier) in your home directory for ASCII rendering. (The app expects it there – customize if needed.)
The source Code Pro Font is Provided in Realeses

To Install Go

Visit the Go downloads page and pick the Windows MSI (amd64 or arm64 depending on your CPU): [https://go.dev/dl/]

Download the .msi and double-click it, follow the installer prompts. After install, restart PowerShell. 

```powershell
   go version
```
Expected Output: go version go1.xx windows/amd64

```powershell
# 1) clone
https://github.com/KrishVij/Clip2ASCII.git
cd .\Clip2ASCII\

# 2) ensure Go modules are synced and dependencies fetched
go mod tidy

# 3) build a Windows executable
go build -o clip2ascii.exe

# 4) run the built executable
.\clip2ascii.exe
```
Notes:

go mod tidy will pull modules like github.com/hajimehoshi/ebiten etc. (no manual git submodules needed).

If your project needs ffmpeg or other system tools, make sure those are installed and available on PATH separately (not covered here).

🐛 Troubleshooting

FFmpeg Not Found: Ensure it's in your PATH. Test with ffmpeg -version.
Font Not Found: Copy a TTF font to ~/Font.ttf.
Video Too Long: App notifies you – trim your video first.
Errors? Check console logs or open an issue with details.

🚀 Usage

Launch the app – you'll see a simple UI with a "Select Video" button.
Choose a video file (formats: MP4, AVI, MKV, etc.). The app checks if it's ≤30 seconds.
Preview the thumbnail in the UI.
Click "ToASCII" to start conversion.

Progress bar shows initial loading.
Real-time updates like "Frame: X processed successfully".


Once done, find your ASCII video in ~/your_video_name.mp4 (home directory).
Close the app – temporary files are auto-deleted.

Pro Tip: For best results, use high-contrast, colorful videos. The ASCII charset is optimized for brightness gradients!

🔧 Under the hood — Detailed (showing the effort)

Below is a slightly expanded, technical-but-readable “behind the scenes” you can paste into your README to show the thought and engineering work that went into Clip2ASCII.

 UI → action — user flow & safety

The Ebiten UI presents a simple flow: pick input, choose output options (font, scale, character set, target FPS), then Convert.

UI code validates inputs, creates well-named temporary folders, and spawns the conversion pipeline while keeping the event loop responsive. Progress, ETA, and per-stage status are reported back to the screen so users always know what’s happening. 

 Extract frames — reliable, reproducible ffmpeg calls

The app shells out to ffmpeg (and ffprobe when needed) to extract a thumbnail and sequential PNG frames into a temp directory. Commands are carefully chosen for reproducibility (explicit frame rate, zero-compression PNGs when desired, predictable file names). 

We keep a thumbnail (quick preview) to update the UI instantly while the heavy work runs in the background. Temporary files are isolated per-run so multiple conversions won’t collide. 

Frame processing (concurrent) — correctness + performance

Frames are processed concurrently with a controlled worker pool (goroutines + semaphore / sync.WaitGroup) so we saturate CPU cores without exhausting memory or I/O. The worker count defaults to runtime.NumCPU() but can be tuned. 

Each worker:

Loads a PNG into memory.

Resizes it to match the target ASCII grid while correcting for character aspect ratio (characters are taller than they are wide). We use a high-quality resampling (Lanczos) to preserve visual detail. 

Subdivides the resized image into blocks (one block per character) and computes a brightness/density value and average color per block. The brightness mapping includes gamma correction and optional contrast tweaks to produce better-looking ASCII density. 

ASCII mapping & character set design — visual finesse

Brightness → character mapping is configurable: multiple character ramps are supported (e.g., @%#*+=-:. ) so users can choose denser or sparser looks. Mapping uses normalized luminance and can include contrast/gamma adjustments for artistic control. 

We preserve color by sampling the average color inside each cell and applying it when drawing the character, so the final result keeps the original frame’s color palette rather than being monochrome. 

Render to text-image — font metrics & pixel-perfect rendering

Characters are drawn onto an RGBA canvas using a TTF font. The renderer uses font metrics (advance, ascender/descender) to align glyphs precisely so blocks line up frame-to-frame and avoid jitter. Font size is chosen to match one canvas cell per character cell, producing consistent alignment. 

We render text with subpixel-aware placement, and then export each rendered frame as PNG ready for re-encoding. This step ensures the ASCII frames are images (so any further video encoder can handle them consistently). 

Stitch & cleanup — final encode and housekeeping

After all ASCII frames are written, the app calls ffmpeg to stitch them back into a video. The stitch command preserves frame rate, optionally sets codec/CRF to control size/quality, and can embed audio from the original if requested. 

All temporary data is removed when the run completes (or when the user cancels), and errors are surfaced to the UI with helpful suggestions (missing ffmpeg on PATH, out-of-disk-space, etc.). 

Concurrency, resource management & UX tradeoffs

The pipeline is designed to balance CPU, memory, and disk I/O. Key tradeoffs:

Higher worker count → faster CPU usage, more memory and I/O pressure.

Larger character cells → fewer glyphs to draw → faster, but lower detail.

Lossless PNG intermediate → larger disk usage but avoids quality loss; user can opt for lossy intermediates in advanced modes.

Error handling & robustness

The app checks ffmpeg/ffprobe exit codes, validates frame sequence integrity before stitching, and implements retries for transient I/O errors. All stages log diagnostic info to a debug log to make reproducing and fixing issues straightforward.

✅ Simple working of the app in flowchart

Extract → resize (aspect-correct) → compute brightness + color → map to ASCII → render with proper font metrics → save PNG frames → stitch with ffmpeg → clean up.

Where to look in the code:

UI & user flow: ui_New.go. 

Frame extraction, ffmpeg calls, and cleanup: FFmpeg_Operations.go. 

Image resizing, brightness → char mapping, and rendering: Processing.go.






