<script lang="ts">
    import { onMount } from "svelte";
    import { ArrowLeft } from "@lucide/svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import Footer from "../components/Footer.svelte";
    import ToastContainer from "../components/ToastContainer.svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import { addToast } from "../stores/toast-wrapper";
    import {
        updateProfile,
        changePassword,
        deactivateAccount,
        updateSettings,
        getProfile,
        getSelfProfile,
        getSettings,
    } from "../controllers/user-controller";
    import type { Settings } from "../types/user";

    let activeTab: string = "profile";
    let theme: "light" | "dark" = "dark";
    let isLoading: boolean = false;

    let profile = {
        name: "",
        username: "",
        bio: "",
        gender: "",
        dob: new Date(),
    };

    let passwordForm = {
        email: "",
        oldPassword: "",
        newPassword: "",
        confirmPassword: "",
    };

    let deactivateForm = {
        password: "",
    };

    let otherSettings: Settings = {
        font_size: "medium",
        font_color: "default",
        private: false,
    };

    async function setActiveTab(tab: string) {
        activeTab = tab;
        isLoading = true;

        try {
            if (tab === "settings") {
                const response = await getSettings();
                if (response.success) {
                    otherSettings = {
                        font_size: response.payload?.font_size || "medium",
                        font_color: response.payload?.font_color || "default",
                        private: response.payload?.private ?? false,
                    };
                } else {
                    addToast("error", response.message, "Settings Fetch Error");
                }
            } else {
                const userId = localStorage.getItem("userId") || "";
                const response = await getSelfProfile(userId);
                if (response.success) {
                    profile = {
                        name: response.payload?.name || "",
                        username: response.payload?.username || "",
                        bio: response.payload?.bio || "",
                        gender: response.payload?.gender || "",
                        dob:
                            new Date(response.payload?.date_of_birth ?? "") ||
                            new Date(),
                    };
                    passwordForm.email = response.payload?.email || "";
                } else {
                    addToast("error", response.message, "Profile Fetch Error");
                }
            }
        } catch (error) {
            addToast("error", "Failed to fetch data", "Error");
            console.error(error);
        } finally {
            isLoading = false;
        }
    }

    async function handleUpdateProfile() {
        isLoading = true;
        try {
            const response = await updateProfile(
                profile.name,
                profile.username,
                profile.bio,
                profile.dob,
                profile.gender
            );
            addToast(
                response.success ? "success" : "error",
                response.success
                    ? "Profile updated successfully!"
                    : response.message,
                response.success ? "Success!" : "Update Failed!"
            );
        } catch (error) {
            addToast("error", "Failed to update profile", "Error");
            console.error(error);
        } finally {
            isLoading = false;
        }
    }

    async function handleChangePassword() {
        if (passwordForm.newPassword !== passwordForm.confirmPassword) {
            addToast("error", "Passwords do not match", "Validation Error");
            return;
        }

        isLoading = true;
        try {
            const response = await changePassword(
                passwordForm.email,
                passwordForm.oldPassword,
                passwordForm.newPassword
            );
            if (response.success) {
                passwordForm = {
                    ...passwordForm,
                    oldPassword: "",
                    newPassword: "",
                    confirmPassword: "",
                };
                addToast(
                    "success",
                    "Password changed successfully!",
                    "Success"
                );
            } else {
                addToast("error", response.message, "Change Failed");
            }
        } catch (error) {
            addToast("error", "Failed to change password", "Error");
            console.error(error);
        } finally {
            isLoading = false;
        }
    }

    async function handleDeactivateAccount() {
        isLoading = true;
        try {
            const response = await deactivateAccount(deactivateForm.password);
            if (response.success) {
                deactivateForm.password = "";
                addToast(
                    "success",
                    "Deactivation request submitted. Check your email.",
                    "Success"
                );
            } else {
                addToast("error", response.message, "Deactivation Failed");
            }
        } catch (error) {
            addToast("error", "Failed to deactivate account", "Error");
            console.error(error);
        } finally {
            isLoading = false;
        }
    }

    async function saveOtherSettings() {
        isLoading = true;
        try {
            const response = await updateSettings(
                otherSettings.font_size,
                otherSettings.font_color,
                otherSettings.private
            );
            if (response.success) {
                addToast("success", "Settings saved successfully!", "Success");
            } else {
                addToast("error", response.message, "Save Failed");
            }
        } catch (error) {
            addToast("error", "Failed to save settings", "Error");
            console.error(error);
        } finally {
            isLoading = false;
        }
    }

    onMount(async () => {
        if (!(await isLoggedIn())) {
            addToast(
                "error",
                "Please log in to access settings",
                "Authentication Error"
            );
            setTimeout(() => (window.location.href = "/login"), 1500);
            return;
        }

        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "dark";
        theme = savedTheme;
        document.documentElement.setAttribute("data-theme", savedTheme);

        await setActiveTab("profile");
    });
</script>

<div class="settings-container">
    <ToastContainer />

    <header class="settings-header">
        <div class="back-button">
            <a href="/home">
                <ArrowLeft size={24} />
            </a>
        </div>
        <h1>Settings</h1>
        <div class="theme-toggle">
            <ToggleTheme bind:theme />
        </div>
    </header>

    <div class="settings-content">
        <div class="tabs">
            <button
                class={activeTab === "profile" ? "active" : ""}
                on:click={() => setActiveTab("profile")}
            >
                Change Profile
            </button>
            <button
                class={activeTab === "password" ? "active" : ""}
                on:click={() => setActiveTab("password")}
            >
                Change Password
            </button>
            <button
                class={activeTab === "deactivate" ? "active" : ""}
                on:click={() => setActiveTab("deactivate")}
            >
                Deactivate Account
            </button>
            <button
                class={activeTab === "settings" ? "active" : ""}
                on:click={() => setActiveTab("settings")}
            >
                Other Settings
            </button>
        </div>

        <div class="tab-content">
            {#if activeTab === "profile"}
                <form on:submit|preventDefault={handleUpdateProfile}>
                    <div class="form-group">
                        <label for="name">Name</label>
                        <input
                            type="text"
                            id="name"
                            bind:value={profile.name}
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="username">Username</label>
                        <input
                            type="text"
                            id="username"
                            bind:value={profile.username}
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="bio">Bio</label>
                        <textarea id="bio" bind:value={profile.bio} rows="4"
                        ></textarea>
                    </div>

                    <div class="form-group">
                        <label for="gender">Gender</label>
                        <select id="gender" bind:value={profile.gender}>
                            <option value="">Select gender</option>
                            <option value="male">Male</option>
                            <option value="female">Female</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="dob">Date of Birth</label>
                        <input type="date" id="dob" bind:value={profile.dob} />
                    </div>

                    <div class="form-actions">
                        <button type="submit" disabled={isLoading}>
                            {isLoading ? "Saving..." : "Save Changes"}
                        </button>
                    </div>
                </form>
            {:else if activeTab === "password"}
                <form on:submit|preventDefault={handleChangePassword}>
                    <div class="form-group">
                        <label for="email">Email</label>
                        <input
                            type="email"
                            id="email"
                            bind:value={passwordForm.email}
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="old-password">Current Password</label>
                        <input
                            type="password"
                            id="old-password"
                            bind:value={passwordForm.oldPassword}
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="new-password">New Password</label>
                        <input
                            type="password"
                            id="new-password"
                            bind:value={passwordForm.newPassword}
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="confirm-password"
                            >Confirm New Password</label
                        >
                        <input
                            type="password"
                            id="confirm-password"
                            bind:value={passwordForm.confirmPassword}
                            required
                        />
                    </div>

                    <div class="form-actions">
                        <button type="submit" disabled={isLoading}>
                            {isLoading ? "Changing..." : "Change Password"}
                        </button>
                    </div>
                </form>
            {:else if activeTab === "deactivate"}
                <div class="deactivate-warning">
                    <p>
                        Warning: Deactivating your account will make your
                        profile and posts inaccessible. This action cannot be
                        undone.
                    </p>
                </div>

                <form on:submit|preventDefault={handleDeactivateAccount}>
                    <div class="form-group">
                        <label for="deactivate-password"
                            >Enter your password to confirm</label
                        >
                        <input
                            type="password"
                            id="deactivate-password"
                            bind:value={deactivateForm.password}
                            required
                        />
                    </div>

                    <div class="form-actions">
                        <button
                            type="submit"
                            class="danger"
                            disabled={isLoading}
                        >
                            {isLoading ? "Processing..." : "Deactivate Account"}
                        </button>
                    </div>
                </form>
            {:else if activeTab === "settings"}
                <form on:submit|preventDefault={saveOtherSettings}>
                    <div class="form-group">
                        <label for="font-size">Font Size</label>
                        <select
                            id="font-size"
                            bind:value={otherSettings.font_size}
                        >
                            <option value="small">Small</option>
                            <option value="medium">Medium</option>
                            <option value="large">Large</option>
                            <option value="x-large">Extra Large</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="font-color">Font Color</label>
                        <select
                            id="font-color"
                            bind:value={otherSettings.font_color}
                        >
                            <option value="black">Black</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="private">Private Account</label>
                        <input
                            type="checkbox"
                            id="private"
                            bind:checked={otherSettings.private}
                        />
                    </div>

                    <div class="form-actions">
                        <button type="submit" disabled={isLoading}>
                            {isLoading ? "Saving..." : "Save Settings"}
                        </button>
                    </div>
                </form>
            {/if}
        </div>
    </div>

    <Footer />
</div>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/settings.scss";
</style>
