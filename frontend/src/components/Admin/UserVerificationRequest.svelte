<script lang="ts">
    import { onMount } from "svelte";
    import { VERIFICATION_IMG } from "../../env_var";
    import type { UserProfile, VerifyAccountRequest } from "../../types/user";
    import { getUserById } from "../../controllers/user-controller";
    import {
        approveUserVerificationRequest,
        rejectUserVerificationRequest,
    } from "../../controllers/admin-controller";
    import { addToast } from "../../stores/toast-wrapper";

    const { request }: { request: VerifyAccountRequest } = $props();
    let premiumRejectionReasons: Record<string, string> = {};

    let user: UserProfile | null = $state(null);

    const handleAccept = async (id: string) => {
        const response = await approveUserVerificationRequest(id);

        if (response.success) {
            addToast("success", "Request accepted successfully!");
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } else {
            addToast("error", response.message, "Error accepting request!");
        }
    };
    const handleReject = async (id: string) => {
        const reason = premiumRejectionReasons[id];
        if (!reason || reason.trim() === "") {
            alert("Please provide a rejection reason.");
            return;
        }
        const response = await rejectUserVerificationRequest(id, reason);

        if (response.success) {
            addToast("success", "Request rejected successfully!");
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } else {
            addToast("error", response.message, "Error rejecting request!");
        }
    };

    onMount(async () => {
        const response = await getUserById(request.user_id);
        if (response) {
            user = response;
        } else {
            user = null;
        }
    });
</script>

<div class="request-card">
    <div class="request-header">
        <div>
            <strong>Request ID - </strong>
            {request.id}
        </div>

        <div class="timestamp">
            {new Date(request.submitted_at).toLocaleString()}
        </div>
    </div>

    <div class="request-info">
        <div>
            User: <a href={"/profile/" + user?.username} target="_blank"
                >{user?.username}</a
            >
        </div>

        <div><strong>Status:</strong> {request.status}</div>
        <div>
            <strong>Reason:</strong>
            {request.reason_text}
        </div>
        <img
            class="selfie-preview"
            src={VERIFICATION_IMG + request.selfie_url}
            alt="Selfie"
        />
    </div>

    <div class="actions">
        <button class="accept-btn" onclick={() => handleAccept(request.id)}
            >Accept</button
        >

        <textarea
            placeholder="Enter rejection reason..."
            bind:value={premiumRejectionReasons[request.id]}
        ></textarea>

        <button class="reject-btn" onclick={() => handleReject(request.id)}
            >Reject</button
        >
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../../styles/admin.scss";
</style>
