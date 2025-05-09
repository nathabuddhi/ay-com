<script lang="ts">
    import { onMount } from "svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import "../styles/app.scss";
    import {
        requestVerificationCode,
        validateVerificationCode,
    } from "../controllers/user-controller";
    import ToastContainer from "../components/ToastContainer.svelte";
    import { addToast } from "../stores/toast-wrapper";

    let currentTheme: "light" | "dark" = "dark";
    let email = "";
    let code = ["", "", "", "", "", ""];
    let isSubmitting = false;
    let errorMessage = "";
    let successMessage = "";
    let timeLeft = 0;
    let timerInterval: number;

    onMount(async () => {
        if (await isLoggedIn()) {
            window.location.href = "/home";
        }

        email = sessionStorage.getItem("verificationEmail") || "";
        if (!email) {
            addToast(
                "error",
                "Email not found. Please register first.",
                "Error"
            );
            setTimeout(() => {
                window.location.href = "/register";
            }, 2000);
            return;
        }

        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "light";
        currentTheme = savedTheme;
        document.documentElement.setAttribute("data-theme", savedTheme);

        startTimer(60);
        addToast(
            "info",
            `We've sent a verification code to ${email}`,
            "Check Your Email"
        );
    });

    function startTimer(seconds: number): void {
        timeLeft = seconds;
        clearInterval(timerInterval);

        timerInterval = setInterval(() => {
            timeLeft--;
            if (timeLeft <= 0) {
                clearInterval(timerInterval);
            }
        }, 1000) as unknown as number;
    }

    function handleInputChange(index: number, event: Event): void {
        const target = event.target as HTMLInputElement;
        const value = target.value;

        if (!/^\d*$/.test(value)) {
            code[index] = "";
            return;
        }

        if (value && index < 5) {
            const nextInput = document.getElementById(
                `code-${index + 1}`
            ) as HTMLInputElement;
            if (nextInput) {
                nextInput.focus();
            }
        }
    }

    function handleKeyDown(index: number, event: KeyboardEvent): void {
        if (event.key === "Backspace" && !code[index] && index > 0) {
            const prevInput = document.getElementById(
                `code-${index - 1}`
            ) as HTMLInputElement;
            if (prevInput) {
                prevInput.focus();
            }
        }
    }

    function handlePaste(event: ClipboardEvent): void {
        event.preventDefault();
        const pastedData = event.clipboardData?.getData("text");

        if (pastedData && /^\d{6}$/.test(pastedData)) {
            for (let i = 0; i < 6; i++) {
                code[i] = pastedData[i];
            }
        }
    }

    async function verifyCode(): Promise<void> {
        const fullCode = code.join("");

        if (fullCode.length !== 6) {
            errorMessage = "Please enter a valid 6-digit code";
            return;
        }

        isSubmitting = true;
        errorMessage = "";
        successMessage = "";

        try {
            const response = await validateVerificationCode(email, fullCode);

            if (response.success) {
                successMessage =
                    "Verification successful! Redirecting to login...";

                setTimeout(() => {
                    window.location.href = "/login";
                }, 2000);
            } else {
                addToast(
                    "error",
                    response.message || "Invalid verification code",
                    "Verification Failed"
                );
            }
        } catch (error) {
            console.error("Verification error:", error);
            addToast(
                "error",
                "An error occurred during verification. Please try again.",
                "Error"
            );
        } finally {
            isSubmitting = false;
        }
    }

    async function requestNewCode(): Promise<void> {
        if (timeLeft > 0) return;

        isSubmitting = true;
        errorMessage = "";
        successMessage = "";

        try {
            const response = await requestVerificationCode(email);

            if (response.success) {
                successMessage =
                    "A new verification code has been sent to your email";
                startTimer(60);
            } else {
                errorMessage =
                    response.message || "Failed to send verification code";
            }
        } catch (error) {
            console.error("Request code error:", error);
            errorMessage = "An error occurred. Please try again.";
        } finally {
            isSubmitting = false;
        }
    }
</script>

<div class="verification-container">
    <ToastContainer />
    <header class="verification-header">
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

    <main class="verification-content">
        <div class="verification-form">
            <div>
                <h1 class="form-title">Verify your email</h1>
                <p class="form-subtitle">
                    We've sent a 6-digit verification code to <strong
                        >{email}</strong
                    >
                </p>
            </div>

            <div class="code-input-container" on:paste={handlePaste}>
                {#each Array(6) as _, i}
                    <input
                        type="text"
                        id="code-{i}"
                        class="code-input"
                        maxlength="1"
                        bind:value={code[i]}
                        on:input={(e) => handleInputChange(i, e)}
                        on:keydown={(e) => handleKeyDown(i, e)}
                        autocomplete="off"
                    />
                {/each}
            </div>

            {#if errorMessage}
                <p class="error-message">{errorMessage}</p>
            {/if}

            {#if successMessage}
                <p class="success-message">{successMessage}</p>
            {/if}

            <button
                class="verify-button"
                on:click={verifyCode}
                disabled={isSubmitting || code.join("").length !== 6}
            >
                {isSubmitting ? "Verifying..." : "Verify"}
            </button>

            <div class="resend-container">
                <p class="resend-text">Didn't receive the code?</p>
                {#if timeLeft > 0}
                    <p class="timer">Resend code in {timeLeft} seconds</p>
                {:else}
                    <button
                        class="resend-button"
                        on:click={requestNewCode}
                        disabled={isSubmitting || timeLeft > 0}
                    >
                        Resend verification code
                    </button>
                {/if}
            </div>
        </div>
    </main>
</div>

<style lang="scss">
    @use "../styles/verification.scss";
</style>
