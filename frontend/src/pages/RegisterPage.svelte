<script lang="ts">
    import { onMount } from "svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import "../styles/app.scss";
    import { register } from "../controllers/user-controller";
    import ToastContainer from "../components/ToastContainer.svelte";
    import { addToast } from "../stores/toast-wrapper";
    import { theme } from "../stores/theme-wrapper";

    let email = $state("");
    let fullName = $state("");
    let username = $state("");
    let password = $state("");
    let confirmPassword = $state("");
    let gender = $state("male");
    let dateOfBirth = $state("");
    let securityQuestion = $state("");
    let securityAnswer = $state("");
    let avatar: File | null = $state(null);
    let banner: File | null = $state(null);

    let errors = $state({
        email: "",
        fullName: "",
        username: "",
        password: "",
        confirmPassword: "",
        dateOfBirth: "",
        securityQuestion: "",
        securityAnswer: "",
        avatar: "",
        banner: "",
    });

    let isSubmitting = $state(false);
    let avatarPreview = $state("");
    let bannerPreview = $state("");

    const securityQuestions = [
        "What was the name of your first pet?",
        "What city were you born in?",
        "What is your favorite video game?",
        "What was the name of your first school?",
        "What was your childhood nickname?",
    ];

    onMount(async () => {
        if (await isLoggedIn()) {
            window.location.href = "/home";
        }

        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "light";
        document.documentElement.setAttribute("data-theme", savedTheme);
    });

    function validateForm(): boolean {
        errors = {
            email: "",
            fullName: "",
            username: "",
            password: "",
            confirmPassword: "",
            dateOfBirth: "",
            securityQuestion: "",
            securityAnswer: "",
            avatar: "",
            banner: "",
        };

        let valid = true;
        let errorMessages: string[] = [];

        if (!email) {
            errors.email = "Email is required";
            errorMessages.push("Email is required");
            valid = false;
        } else if (!/\S+@\S+\.\S+/.test(email)) {
            errors.email = "Email is invalid";
            errorMessages.push("Email is invalid");
            valid = false;
        }

        if (!fullName) {
            errors.fullName = "Full name is required";
            errorMessages.push("Full name is required");
            valid = false;
        }

        if (!username) {
            errors.username = "Username is required";
            errorMessages.push("Username is required");
            valid = false;
        } else if (username.length < 3) {
            errors.username = "Username must be at least 3 characters";
            errorMessages.push("Username must be at least 3 characters");
            valid = false;
        }

        if (!password) {
            errors.password = "Password is required";
            errorMessages.push("Password is required");
            valid = false;
        } else if (password.length < 8) {
            errors.password = "Password must be at least 8 characters";
            errorMessages.push("Password must be at least 8 characters");
            valid = false;
        }

        if (password !== confirmPassword) {
            errors.confirmPassword = "Passwords do not match";
            errorMessages.push("Passwords do not match");
            valid = false;
        }

        if (!dateOfBirth) {
            errors.dateOfBirth = "Date of birth is required";
            errorMessages.push("Date of birth is required");
            valid = false;
        } else {
            const dob = new Date(dateOfBirth);
            const today = new Date();
            const age = today.getFullYear() - dob.getFullYear();
            if (age < 13) {
                errors.dateOfBirth = "You must be at least 13 years old";
                errorMessages.push("You must be at least 13 years old");
                valid = false;
            }
        }

        if (!securityQuestion) {
            errors.securityQuestion = "Security question is required";
            errorMessages.push("Security question is required");
            valid = false;
        }

        if (!securityAnswer) {
            errors.securityAnswer = "Security answer is required";
            errorMessages.push("Security answer is required");
            valid = false;
        }

        if (!avatar) {
            errors.avatar = "Profile picture is required";
            errorMessages.push("Profile picture is required");
            valid = false;
        }

        if (!banner) {
            errors.banner = "Banner image is required";
            errorMessages.push("Banner image is required");
            valid = false;
        }

        return valid;
    }

    function handleAvatarChange(event: Event): void {
        const target = event.target as HTMLInputElement;
        if (target.files && target.files[0]) {
            avatar = target.files[0];

            const reader = new FileReader();
            reader.onload = (e) => {
                avatarPreview = e.target?.result as string;
            };
            reader.readAsDataURL(avatar);
            errors.avatar = "";
        }
    }

    function handleBannerChange(event: Event): void {
        const target = event.target as HTMLInputElement;
        if (target.files && target.files[0]) {
            banner = target.files[0];

            const reader = new FileReader();
            reader.onload = (e) => {
                bannerPreview = e.target?.result as string;
            };
            reader.readAsDataURL(banner);
            errors.banner = "";
        }
    }

    async function handleSubmit(): Promise<void> {
        if (!validateForm()) {
            addToast(
                "error",
                "Please fix the errors in the form.",
                "Invalid Input"
            );
            return;
        }

        isSubmitting = true;

        try {
            const response = await register(
                email,
                fullName,
                username,
                password,
                gender,
                new Date(dateOfBirth),
                securityQuestion,
                securityAnswer,
                avatar as File,
                banner as File
            );

            if (response.success) {
                addToast(
                    "success",
                    "Registration successful! Please verify your email.",
                    "Success"
                );

                sessionStorage.setItem("verificationEmail", email);

                setTimeout(() => {
                    window.location.href = "/verification";
                }, 1500);
            } else {
                addToast("error", response.message, "Registration Failed");
            }
        } catch (error) {
            addToast(
                "error",
                "An error occurred during registration. Please try again.",
                "Error"
            );
        } finally {
            isSubmitting = false;
        }
    }
</script>

<div class="register-container">
    <header class="register-header">
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

    <main class="register-content">
        <ToastContainer />

        <h1 class="form-title">Create your account</h1>

        <form class="register-form" onsubmit={handleSubmit}>
            <div class="form-group">
                <label for="email" class="form-label">Email</label>
                <input
                    type="email"
                    id="email"
                    class="form-input"
                    bind:value={email}
                    onblur={validateForm}
                    required
                />
                {#if errors.email}
                    <span class="error-message">{errors.email}</span>
                {/if}
            </div>

            <div class="form-group">
                <label for="fullName" class="form-label">Full Name</label>
                <input
                    type="text"
                    id="fullName"
                    class="form-input"
                    bind:value={fullName}
                    onblur={validateForm}
                    required
                />
                {#if errors.fullName}
                    <span class="error-message">{errors.fullName}</span>
                {/if}
            </div>

            <div class="form-group">
                <label for="username" class="form-label">Username</label>
                <input
                    type="text"
                    id="username"
                    class="form-input"
                    bind:value={username}
                    onblur={validateForm}
                    required
                />
                {#if errors.username}
                    <span class="error-message">{errors.username}</span>
                {/if}
            </div>

            <div class="form-group">
                <label for="password" class="form-label">Password</label>
                <input
                    type="password"
                    id="password"
                    class="form-input"
                    bind:value={password}
                    onblur={validateForm}
                    required
                />
                {#if errors.password}
                    <span class="error-message">{errors.password}</span>
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
                    bind:value={confirmPassword}
                    onblur={validateForm}
                    required
                />
                {#if errors.confirmPassword}
                    <span class="error-message">{errors.confirmPassword}</span>
                {/if}
            </div>

            <div class="form-group">
                <label for="dateOfBirth" class="form-label">Date of Birth</label
                >
                <input
                    type="date"
                    id="dateOfBirth"
                    class="form-input"
                    bind:value={dateOfBirth}
                    onblur={validateForm}
                    required
                />
                {#if errors.dateOfBirth}
                    <span class="error-message">{errors.dateOfBirth}</span>
                {/if}
            </div>

            <div class="form-group">
                <label class="form-label">Gender</label>
                <div class="form-radio-group">
                    <label class="radio-option">
                        <input
                            type="radio"
                            name="gender"
                            value="male"
                            bind:group={gender}
                        />
                        Male
                    </label>
                    <label class="radio-option">
                        <input
                            type="radio"
                            name="gender"
                            value="female"
                            bind:group={gender}
                        />
                        Female
                    </label>
                </div>
            </div>

            <div class="form-group">
                <label for="securityQuestion" class="form-label"
                    >Security Question</label
                >
                <select
                    id="securityQuestion"
                    class="form-select"
                    bind:value={securityQuestion}
                    onblur={validateForm}
                    required
                >
                    <option value="" disabled selected
                        >Select a security question</option
                    >
                    {#each securityQuestions as question}
                        <option value={question}>{question}</option>
                    {/each}
                </select>
                {#if errors.securityQuestion}
                    <span class="error-message">{errors.securityQuestion}</span>
                {/if}
            </div>

            <div class="form-group">
                <label for="securityAnswer" class="form-label"
                    >Security Answer</label
                >
                <input
                    type="text"
                    id="securityAnswer"
                    class="form-input"
                    bind:value={securityAnswer}
                    onblur={validateForm}
                    required
                />
                {#if errors.securityAnswer}
                    <span class="error-message">{errors.securityAnswer}</span>
                {/if}
            </div>

            <div class="form-group">
                <label class="form-label">Profile Picture</label>
                <div class="file-input-container">
                    <label for="avatar" class="file-input-label">
                        {avatar ? avatar.name : "Choose a profile picture"}
                    </label>
                    <input
                        type="file"
                        id="avatar"
                        class="file-input"
                        accept="image/*"
                        onchange={handleAvatarChange}
                    />
                    {#if avatarPreview}
                        <img
                            src={avatarPreview || "/placeholder.svg"}
                            alt="Avatar preview"
                            class="file-preview visible"
                        />
                    {/if}
                </div>
                {#if errors.avatar}
                    <span class="error-message">{errors.avatar}</span>
                {/if}
            </div>

            <div class="form-group">
                <label class="form-label">Banner Image</label>
                <div class="file-input-container">
                    <label for="banner" class="file-input-label">
                        {banner ? banner.name : "Choose a banner image"}
                    </label>
                    <input
                        type="file"
                        id="banner"
                        class="file-input"
                        accept="image/*"
                        onchange={handleBannerChange}
                    />
                    {#if bannerPreview}
                        <img
                            src={bannerPreview || "/placeholder.svg"}
                            alt="Banner preview"
                            class="file-preview visible"
                        />
                    {/if}
                </div>
                {#if errors.banner}
                    <span class="error-message">{errors.banner}</span>
                {/if}
            </div>

            <button
                type="submit"
                class="register-button"
                disabled={isSubmitting}
            >
                {isSubmitting ? "Creating account..." : "Create account"}
            </button>
        </form>

        <div class="login-link">
            Already have an account? <a href="/login">Sign in</a>
        </div>
        <div class="login-link">
            <a href="/">Back to Landing Page</a>
        </div>
    </main>
</div>

<style lang="scss">
    @use "../styles/register.scss";
</style>
