<script lang="ts">
    import { CircleStop, Mic } from "@lucide/svelte";
    import { addToast } from "../stores/toast-wrapper";

    let mediaRecorder: MediaRecorder;
    let chunks: BlobPart[] = [];
    let isRecording = $state(false);
    let loading = $state(false);

    let { query = $bindable() }: { query: string } = $props();

    async function startRecording() {
        const stream = await navigator.mediaDevices.getUserMedia({
            audio: true,
        });
        mediaRecorder = new MediaRecorder(stream);

        mediaRecorder.ondataavailable = (e) => chunks.push(e.data);

        mediaRecorder.onstop = async () => {
            const blob = new Blob(chunks, { type: "audio/webm" });
            chunks = [];
            const file = new File([blob], "recording.webm", {
                type: "audio/webm",
            });

            const formData = new FormData();
            formData.append("audio", file);

            loading = true;
            try {
                const res = await fetch("http://localhost:5020/transcribe", {
                    method: "POST",
                    body: formData,
                });
                const data = await res.json();
                query = data.transcription || "";
                query = query.trim().toLowerCase();
            } catch (err) {
                addToast("error", "Transcription failed. Please Try Again.");
            }
            loading = false;
        };

        mediaRecorder.start();
        isRecording = true;
    }

    function stopRecording() {
        if (mediaRecorder && mediaRecorder.state !== "inactive") {
            mediaRecorder.stop();
            isRecording = false;
            addToast("info", "Recording stopped. Transcribing...");
        }
    }
</script>

{#if isRecording}
    <button onclick={stopRecording} disabled={!isRecording} class="icon-button">
        <CircleStop />
    </button>
{:else}
    <button onclick={startRecording} disabled={isRecording} class="icon-button"
        ><Mic /></button
    >
{/if}

<style>
    .icon-button {
        background: var(--primary-hover);
        border: none;
        cursor: pointer;
        padding: 0.5rem;
        display: flex;
        border-radius: 100%;
        margin: 0.5rem;
    }
</style>
