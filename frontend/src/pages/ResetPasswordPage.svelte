<script lang="ts">
    import { onMount } from "svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import Footer from "../components/Footer.svelte";
    import ToastContainer from "../components/ToastContainer.svelte";
    import { addToast } from "../stores/toast-wrapper";
    import { isLoggedIn } from "../controllers/token-controller";
    import "../styles/app.scss";
    import { resetPassword } from "../controllers/user-controller";

    let currentTheme: "light" | "dark" = "dark";
    let email = "";
    let code = "";
    let password = "";
    let confirmPassword = "";
    let isSubmitting = false;
    let passwordStrength = "";

    let errors = {
        email: "",
        code: "",
        password: "",
        confirmPassword: "",
    };

    onMount(async () => {
        if (await isLoggedIn()) {
            window.location.href = "/home";
        }

        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "light";
        currentTheme = savedTheme;
        document.documentElement.setAttribute("data-theme", savedTheme);

        const urlParams = new URLSearchParams(window.location.search);
        const emailParam = urlParams.get("email");
        if (emailParam) {
            email = emailParam;
        }

        const codeParam = urlParams.get("code");
        if (codeParam) {
            code = codeParam;
        }
    });

    function validateForm(): boolean {
        errors = {
            email: "",
            code: "",
            password: "",
            confirmPassword: "",
        };

        let valid = true;

        if (!email) {
            errors.email = "Email is required";
            valid = false;
        } else if (!/\S+@\S+\.\S+/.test(email)) {
            errors.email = "Email is invalid";
            valid = false;
        }

        if (!code) {
            errors.code = "Reset code is required";
            valid = false;
        } else if (code.length !== 10 || !/^\d+$/.test(code)) {
            errors.code = "Reset code must be 10 digits";
            valid = false;
        }

        if (!password) {
            errors.password = "Password is required";
            valid = false;
        } else if (password.length < 8) {
            errors.password = "Password must be at least 8 characters";
            valid = false;
        }

        if (password !== confirmPassword) {
            errors.confirmPassword = "Passwords do not match";
            valid = false;
        }

        return valid;
    }

    function checkPasswordStrength(password: string): void {
        if (!password) {
            passwordStrength = "";
            return;
        }

        const hasLowerCase = /[a-z]/.test(password);
        const hasUpperCase = /[A-Z]/.test(password);
        const hasNumbers = /\d/.test(password);
        const hasSpecialChars = /[!@#$%^&*(),.?":{}|<>]/.test(password);

        const strength =
            (hasLowerCase ? 1 : 0) +
            (hasUpperCase ? 1 : 0) +
            (hasNumbers ? 1 : 0) +
            (hasSpecialChars ? 1 : 0);

        if (password.length < 8 || strength <= 1) {
            passwordStrength = "weak";
        } else if (strength === 2) {
            passwordStrength = "medium";
        } else {
            passwordStrength = "strong";
        }
    }

    async function handleSubmit(): Promise<void> {
        if (!validateForm()) {
            addToast(
                "error",
                "Please fix the errors in the form",
                "Validation Error"
            );
            return;
        }

        isSubmitting = true;

        try {
            const response = await resetPassword(email, code, password);

            if (response.success) {
                addToast(
                    "success",
                    "Password reset successful! You can now log in with your new password.",
                    "Success"
                );

                setTimeout(() => {
                    window.location.href = "/login";
                }, 2000);
            } else {
                addToast(
                    "error",
                    response.message || "Failed to reset password",
                    "Error"
                );
            }
        } catch (error) {
            addToast(
                "error",
                error instanceof Error
                    ? error.message
                    : "An unknown error occurred.",
                "Error"
            );
        } finally {
            isSubmitting = false;
        }
    }
</script>

<div class="reset-password-container">
    <ToastContainer />

    <header class="reset-password-header">
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

    <main class="reset-password-content">
        <div class="reset-password-form">
            <div>
                <h1 class="form-title">Reset Password</h1>
                <p class="form-subtitle">
                    Enter your email, the reset code from your email, and your
                    new password.
                </p>
            </div>

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
                {#if errors.email}
                    <span class="field-error">{errors.email}</span>
                {/if}
            </div>

            <div class="form-group">
                <label for="code" class="form-label">Reset Code</label>
                <input
                    type="text"
                    id="code"
                    class="form-input code-input"
                    placeholder="Enter the 10-digit code from your email"
                    bind:value={code}
                    maxlength="10"
                    disabled={isSubmitting}
                    required
                />
                {#if errors.code}
                    <span class="field-error">{errors.code}</span>
                {/if}
            </div>

            <div class="form-group">
                <label for="password" class="form-label">New Password</label>
                <input
                    type="password"
                    id="password"
                    class="form-input"
                    placeholder="Enter your new password"
                    bind:value={password}
                    on:input={() => checkPasswordStrength(password)}
                    disabled={isSubmitting}
                    required
                />
                {#if passwordStrength}
                    <div class="password-strength">
                        <span>Password strength: {passwordStrength}</span>
                        <div class="strength-bar">
                            <div
                                class="strength-indicator {passwordStrength}"
                            ></div>
                        </div>
                    </div>
                {/if}
                {#if errors.password}
                    <span class="field-error">{errors.password}</span>
                {/if}
            </div>

            <div class="form-group">
                <label for="confirmPassword" class="form-label"
                    >Confirm Password</label
                >
                <input
                    type="password"
                    id="confirmPassword"
                    class="form-input"
                    placeholder="Confirm your new password"
                    bind:value={confirmPassword}
                    disabled={isSubmitting}
                    required
                />
                {#if errors.confirmPassword}
                    <span class="field-error">{errors.confirmPassword}</span>
                {/if}
            </div>

            <button
                class="submit-button"
                on:click={handleSubmit}
                disabled={isSubmitting}
            >
                {isSubmitting ? "Resetting..." : "Reset Password"}
            </button>
        </div>

        <div class="back-link">
            <a href="/login">Back to Login</a>
        </div>
    </main>

    <Footer />
</div>

<style lang="scss">
    @use "../styles/reset-password.scss";
</style>
