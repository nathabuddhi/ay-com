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

    let currentTheme: "light" | "dark" = "dark";
    let email = "";
    let securityAnswer = "";
    let isSubmitting = false;
    let showSuccessMessage = false;

    let currentStep = 1;

    let securityQuestion = "";

    onMount(async () => {
        if (await isLoggedIn()) {
            window.location.href = "/home";
        }

        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "light";
        currentTheme = savedTheme;
        document.documentElement.setAttribute("data-theme", savedTheme);
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
                src={currentTheme === "dark"
                    ? "logo-dark.png"
                    : "logo-light.png"}
                alt="AY Logo"
                class="logo-img"
            />
        </div>
        <div class="theme-toggle-container">
            <ToggleTheme bind:theme={currentTheme} />
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
                    on:click={handleEmailSubmit}
                    disabled={isSubmitting}
                >
                    {isSubmitting ? "Checking..." : "Continue"}
                </button>
            {:else if currentStep === 2}
                <div class="security-question-container">
                    <div class="security-question">
                        <span class="question-label">Security Question:</span>
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
                        on:click={() => (currentStep = 1)}
                        disabled={isSubmitting}
                    >
                        Back
                    </button>
                    <button
                        class="submit-button"
                        on:click={handleSecurityAnswerSubmit}
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

    .security-question-container {
        display: flex;
        flex-direction: column;
        gap: 16px;
    }

    .security-question {
        background-color: rgba(var(--background-rgb), 0.5);
        border-radius: 8px;
        padding: 16px;
        border: 1px solid var(--border);
    }

    .question-label {
        font-size: 14px;
        font-weight: 500;
        color: var(--secondary);
        display: block;
        margin-bottom: 8px;
    }

    .question-text {
        font-size: 16px;
        font-weight: 500;
        color: var(--text);
    }

    .button-group {
        display: flex;
        gap: 12px;
        margin-top: 8px;
    }

    .back-button {
        padding: 12px 16px;
        border-radius: 9999px;
        background-color: transparent;
        color: var(--text);
        font-size: 16px;
        font-weight: 500;
        border: 1px solid var(--border);
        cursor: pointer;
        transition: background-color 0.2s ease;
        flex: 1;

        &:hover {
            background-color: rgba(var(--background-rgb), 0.1);
        }

        &:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
    }

    .submit-button {
        flex: 2;
    }

    @media (max-width: 768px) {
        .button-group {
            flex-direction: column;
        }

        .back-button,
        .submit-button {
            flex: auto;
        }
    }

    @media (max-width: 480px) {
        .security-question {
            padding: 12px;
        }

        .question-text {
            font-size: 14px;
        }
    }
</style>
