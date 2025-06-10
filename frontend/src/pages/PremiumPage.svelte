<script lang="ts">
    import { onMount } from "svelte";
    import { getValidToken, isLoggedIn } from "../controllers/token-controller";
    import { addToast } from "../stores/toast-wrapper";
    let isVerified = $state(false);
    let idCard = $state("");
    let reason = $state("");
    let faceFile: File | null = null;
    let submitting = $state(false);

    let facePreviewUrl = $state("");

    onMount(async () => {
        if (await !isLoggedIn()) {
            window.location.href = "/login";
        }
        isVerified = localStorage.getItem("is_verified") === "true";
    });

    function handleFileChange(e: Event) {
        faceFile = (e.target as HTMLInputElement).files?.[0] ?? null;
        facePreviewUrl = faceFile ? URL.createObjectURL(faceFile) : "";
    }

    async function submitForm(event: Event) {
        event?.preventDefault();

        if (!idCard || !reason || !faceFile) {
            addToast(
                "error",
                "Please fill in all fields and upload a face photo.",
                "Empty Fields!"
            );
            return;
        }
        if (idCard.length < 7 || idCard.length > 20) {
            addToast(
                "error",
                "ID card number must be between 7 and 20 characters.",
                "Invalid ID!"
            );
            return;
        }

        submitting = true;

        const fd = new FormData();
        fd.append("identity_card_number", idCard);
        fd.append("reason", reason);
        fd.append("face", faceFile);

        try {
            await sendVerificationRequest(fd);
            addToast(
                "success",
                "Verification request submitted successfully!",
                "Success!"
            );
        } catch (error) {
            addToast(
                "error",
                "Failed to submit verification request: " +
                    (error as Error).message,
                "Error"
            );
        } finally {
            submitting = false;
            idCard = "";
            reason = "";
            faceFile = null;
        }
    }

    async function sendVerificationRequest(fd: FormData): Promise<void> {
        try {
            const response = await fetch(
                "http://localhost:5000/user/submitverifyaccountrequest",
                {
                    method: "POST",
                    body: fd,
                    headers: {
                        Authorization: await getValidToken(),
                    },
                }
            );

            const data = await response.json();
            if (!data.success) {
                throw new Error(data.message || "Verification failed");
            }
        } catch (error) {
            throw error;
        }
    }
</script>

{#if isVerified}
    <div class="verified-msg">
        You're already a premium member with a blue checkmark!
    </div>
{:else}
    <div class="form-container">
        <h2>Get Verified & Earn a Blue Checkmark</h2>
        <form onsubmit={submitForm}>
            <label>
                National ID Card Number
                <input type="text" bind:value={idCard} required />
            </label>

            <label>
                Reason for Verification
                <textarea bind:value={reason} rows="3" required></textarea>
            </label>

            <label class="file-upload">
                Upload Face Photo
                <input
                    type="file"
                    accept="image/*"
                    onchange={handleFileChange}
                    required
                />
            </label>
            <img src={facePreviewUrl} alt="Face Preview" class="face-preview" />

            <button type="submit" disabled={submitting}>
                {submitting ? "Submitting..." : "Submit"}
            </button>
        </form>
    </div>
{/if}

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @import "../styles/premium.scss";
</style>
