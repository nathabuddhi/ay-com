<script lang="ts">
    import { onMount } from "svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import ToastContainer from "../components/ToastContainer.svelte";
    import { addToast } from "../stores/toast-wrapper";
    import { isLoggedIn } from "../controllers/token-controller";
    import "../styles/app.scss";
    import {
        getSecurityQuestion,
        validateSecurityQuestion,
    } from "../controllers/user-controller";
    import { theme } from "../stores/theme-wrapper";

    let email = $state("");
    let securityAnswer = $state("");
    let isSubmitting = $state(false);
    let showSuccessMessage = $state(false);

    let currentStep = $state(1);

    let securityQuestion = $state("");

    onMount(async () => {
        if (await isLoggedIn()) {
            window.location.href = "/home";
        }
    });

    function validateEmail(email: string): boolean {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }

    async function handleEmailSubmit(): Promise<void> {
        if (!email) {
            addToast("error", "Please enter your email address", "Error");
            return;
        }

        if (!validateEmail(email)) {
            addToast("error", "Please enter a valid email address", "Error");
            return;
        }

        isSubmitting = true;

        try {
            const response = await getSecurityQuestion(email);

            if (response.success && response.payload) {
                securityQuestion = response.payload.value;
                currentStep = 2;
                addToast(
                    "info",
                    "Please answer your security question",
                    "Security Verification"
                );
            } else {
                addToast(
                    "error",
                    response.message || "Email not found",
                    "Error"
                );
            }
        } catch (error) {
            console.error("Get security question error:", error);
            addToast("error", "An error occurred. Please try again.", "Error");
        } finally {
            isSubmitting = false;
        }
    }

    async function handleSecurityAnswerSubmit(): Promise<void> {
        if (!securityAnswer) {
            addToast("error", "Please enter your security answer", "Error");
            return;
        }

        isSubmitting = true;

        try {
            const response = await validateSecurityQuestion(
                email,
                securityAnswer
            );

            if (response.success) {
                currentStep = 3;
                showSuccessMessage = true;
                addToast(
                    "success",
                    "Password reset instructions sent to your email",
                    "Success"
                );
            } else {
                addToast(
                    "error",
                    response.message || "Incorrect security answer",
                    "Error"
                );
            }
        } catch (error) {
            console.error("Security answer verification error:", error);
            addToast("error", "An error occurred. Please try again.", "Error");
        } finally {
            isSubmitting = false;
        }
    }
</script>

<div class="forgot-password-container">
    <ToastContainer />

    <header class="forgot-password-header">
        <div class="logo-container">
            <img
                src={$theme === "dark" ? "logo-dark.png" : "logo-light.png"}
                alt="AY Logo"
                class="logo-img"
            />
        </div>
        <div class="theme-toggle-container">
            <ToggleTheme />
        </div>
    </header>

    <main class="forgot-password-content">
        <div class="forgot-password-form">
            <div>
                <h1 class="form-title">Forgot Password</h1>
                {#if currentStep === 1}
                    <p class="form-subtitle">
                        Enter your email address to start the password recovery
                        process.
                    </p>
                {:else if currentStep === 2}
                    <p class="form-subtitle">
                        Please answer your security question to verify your
                        identity.
                    </p>
                {:else}
                    <p class="form-subtitle">
                        Password reset instructions have been sent to your
                        email.
                    </p>
                {/if}
            </div>

            {#if currentStep === 1}
                <div class="form-group">
                    <label for="email" class="form-label">Email</label>
                    <input
                        type="email"
                        id="email"
                        class="form-input"
                        placeholder="Enter your email address"
                        bind:value={email}
                        disabled={isSubmitting}
                        required
                    />
                </div>

                <button
                    class="submit-button"
                    onclick={handleEmailSubmit}
                    disabled={isSubmitting}
                >
                    {isSubmitting ? "Checking..." : "Continue"}
                </button>
            {:else if currentStep === 2}
                <div class="security-question-container">
                    <div class="security-question">
                        <span class="question-label">Security Question</span>
                        <p class="question-text">{securityQuestion}</p>
                    </div>

                    <div class="form-group">
                        <label for="securityAnswer" class="form-label"
                            >Your Answer</label
                        >
                        <input
                            type="text"
                            id="securityAnswer"
                            class="form-input"
                            placeholder="Enter your answer"
                            bind:value={securityAnswer}
                            disabled={isSubmitting}
                            required
                        />
                    </div>
                </div>

                <div class="button-group">
                    <button
                        class="back-button"
                        onclick={() => (currentStep = 1)}
                        disabled={isSubmitting}
                    >
                        Back
                    </button>
                    <button
                        class="submit-button"
                        onclick={handleSecurityAnswerSubmit}
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? "Verifying..." : "Submit"}
                    </button>
                </div>
            {:else}
                <div class="success-message visible">
                    <p>
                        We've sent password reset instructions to <strong
                            >{email}</strong
                        >. Please check your email and follow the link to reset
                        your password.
                    </p>
                </div>
            {/if}
        </div>

        <div class="back-link">
            <a href="/login">Back to Login</a>
        </div>
    </main>
</div>

<style lang="scss">
    @use "../styles/forgot-password.scss";
</style>
