<script lang="ts">
    import { AlertTriangle, Ban, Check, UserCheck, X } from "@lucide/svelte";
    import type { UserReport } from "../../types/admin";
    import { addToast } from "../../stores/toast-wrapper";
    import {
        approveReport,
        rejectReport,
    } from "../../controllers/admin-controller";

    let { report }: { report: UserReport } = $props();

    async function handleApproveReport() {
        try {
            const response = await approveReport(report.report_id);
            if (response.success) {
                addToast("success", "Report approved and user banned");
                report.status = "approved";
            } else {
                addToast(
                    "error",
                    response.message || "Failed to approve report"
                );
            }
        } catch (error) {
            addToast("error", "Failed to approve report");
        }
    }

    async function handleRejectReport() {
        try {
            const response = await rejectReport(report.report_id);
            if (response.success) {
                addToast("success", "Report rejected");
                report.status = "rejected";
            } else {
                addToast(
                    "error",
                    response.message || "Failed to reject report"
                );
            }
        } catch (error) {
            addToast("error", "Failed to reject report");
        }
    }
</script>

<div class="report-card">
    <div class="report-header">
        <AlertTriangle size={20} class="warning-icon" />
        <h4>User Report</h4>
    </div>
    <div class="report-info">
        <p>
            <strong>Reporter:</strong>
            {report.reporter_id}
        </p>
        <p>
            <strong>Reported User:</strong>
            {report.reported_id}
        </p>
        <p>
            <strong>Reason:</strong>
            {report.reason}
        </p>
    </div>
    <div class="actions">
        {#if report.status === "approved"}
            <p class="status approved">
                <UserCheck size={16} />
                Report Approved
            </p>
        {:else if report.status === "rejected"}
            <p class="status rejected">
                <Ban size={16} />
                Report Rejected
            </p>
        {:else}
            <p class="status pending">
                <AlertTriangle size={16} />
                Pending Review
            </p>
            <button
                class="action-btn approve-btn"
                onclick={handleApproveReport}
            >
                <Check size={16} />
                Approve & Ban User
            </button>
            <button class="action-btn reject-btn" onclick={handleRejectReport}>
                <X size={16} />
                Reject Report
            </button>
        {/if}
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../../styles/admin.scss";
</style>
