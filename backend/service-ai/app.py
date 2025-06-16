from flask import Flask, request, jsonify
from flask_cors import CORS
from transformers import Wav2Vec2ForCTC, Wav2Vec2Processor
import torch
import torchaudio
import os
import tempfile
import subprocess # <-- 1. Import subprocess, pydub is no longer needed

app = Flask(__name__)
CORS(app)

# 2. Define the path to your FFmpeg executable.
#    Using a simple path is recommended (e.g., move ffmpeg to C:\ffmpeg)
FFMPEG_PATH = "./model/ffmpeg.exe" 
# Or use the full path if you prefer, e.g.:
# FFMPEG_PATH = "C:\\APPS\\Installers\\ffmpeg-build-folder\\bin\\ffmpeg.exe"

# --- Sanity check to ensure the FFmpeg path is valid ---
if not os.path.exists(FFMPEG_PATH):
    raise FileNotFoundError(f"FFmpeg executable not found at path: {FFMPEG_PATH}")
# ---------------------------------------------------------


# Load processor and model (or from local path if already saved)
model_dir = "./model/wav2vec2-960h"
if not os.path.exists(model_dir):
    processor = Wav2Vec2Processor.from_pretrained("facebook/wav2vec2-base-960h")
    model = Wav2Vec2ForCTC.from_pretrained("facebook/wav2vec2-base-960h")
    model.save_pretrained(model_dir)
    processor.save_pretrained(model_dir)
else:
    processor = Wav2Vec2Processor.from_pretrained(model_dir)
    model = Wav2Vec2ForCTC.from_pretrained(model_dir)


@app.route("/transcribe", methods=["POST"])
def transcribe_audio():
    if 'audio' not in request.files:
        return jsonify({"error": "No audio file provided"}), 400

    audio_file = request.files['audio']
    
    temp_dir = tempfile.gettempdir()

    # Save uploaded file (likely WebM or OGG) temporarily
    input_path = os.path.join(temp_dir, audio_file.filename)
    audio_file.save(input_path)

    # Define the output path for the converted WAV file
    output_path = os.path.join(temp_dir, "converted.wav")

    try:
        # 3. This is the new conversion logic using subprocess
        #    Build the FFmpeg command as a list of arguments
        command = [
            FFMPEG_PATH,
            '-i', input_path,      # Input file
            '-ar', '16000',        # Set audio sample rate to 16kHz
            '-ac', '1',            # Set audio channels to 1 (mono)
            '-acodec', 'pcm_s16le',# Set audio codec to 16-bit PCM for WAV
            '-y',                  # Overwrite output file if it exists
            output_path
        ]

        # Execute the command
        # check=True will raise an exception if FFmpeg returns a non-zero exit code
        # capture_output=True will hide FFmpeg's console output from your Flask log
        subprocess.run(command, check=True, capture_output=True)

    except subprocess.CalledProcessError as e:
        # If FFmpeg fails, this will be triggered
        # e.stderr will contain the error message from FFmpeg
        error_message = e.stderr.decode()
        return jsonify({"error": f"Audio conversion failed with FFmpeg: {error_message}"}), 500
    except Exception as e:
        # Catch other potential exceptions
        return jsonify({"error": f"An unexpected error occurred: {str(e)}"}), 500

    # Now load the converted WAV with torchaudio
    speech_array, sampling_rate = torchaudio.load(output_path)

    input_values = processor(speech_array.squeeze().numpy(), return_tensors="pt", sampling_rate=16000).input_values
    with torch.no_grad():
        logits = model(input_values).logits
    predicted_ids = torch.argmax(logits, dim=-1)
    transcription = processor.decode(predicted_ids[0])

    # Optional: Clean up the temporary files
    os.remove(input_path)
    os.remove(output_path)

    return jsonify({"transcription": transcription})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5020, debug=True)