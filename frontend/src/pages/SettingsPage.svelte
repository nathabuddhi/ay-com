<script lang="ts">
    import { onMount } from "svelte";
    import { ArrowLeft } from "@lucide/svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import ToastContainer from "../components/ToastContainer.svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import { addToast } from "../stores/toast-wrapper";
    import {
        updateProfile,
        changePassword,
        deactivateAccount,
        updateSettings,
        getSelfProfile,
        getSettings,
    } from "../controllers/user-controller";
    import type { Settings } from "../types/user";
    import BlockedUser from "../components/BlockedUser.svelte";
    import type { NotificationSettings } from "../types/settings";
    import {
        getNotificationSettings,
        updateNotificationSettings,
    } from "../controllers/notification-controller";

    let activeTab: string = $state("profile");
    let isLoading: boolean = $state(false);

    let profile = $state({
        name: "",
        username: "",
        bio: "",
        gender: "",
        dob: new Date(),
    });

    let passwordForm = $state({
        email: "",
        oldPassword: "",
        newPassword: "",
        confirmPassword: "",
    });

    let deactivateForm = $state({
        password: "",
    });

    let otherSettings = $state({
        font_size: "medium",
        font_color: "black",
        private: false,
    });

    let notificationSettings = $state<NotificationSettings>({
        notif_like: true,
        notif_community: true,
        notif_follow: true,
        notif_mention: true,
        notif_newsletter: true,
        notif_repost: true,
    });

    async function setActiveTab(tab: string) {
        activeTab = tab;
        isLoading = true;

        try {
            if (tab === "settings") {
                const response = await getSettings();
                if (response.success) {
                    otherSettings = {
                        font_size: response.payload?.font_size || "medium",
                        font_color: response.payload?.font_color || "black",
                        private: response.payload?.private ?? false,
                    };
                } else {
                    addToast("error", response.message, "Settings Fetch Error");
                }
            } else if (tab === "notification") {
                const response = await getNotificationSettings();
                if (response.success && response.payload) {
                    notificationSettings = response.payload;
                } else {
                    addToast("error", response.message, "Settings Fetch Error");
                }
            } else {
                const response = await getSelfProfile();
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
            addToast("error", "Failed to update profile: " + error, "Error");
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

    async function saveNotificationSettings() {
        isLoading = true;
        try {
            const response =
                await updateNotificationSettings(notificationSettings);
            if (response.success) {
                addToast("success", "Settings saved successfully!", "Success");
            } else {
                addToast("error", response.message, "Save Failed");
            }
        } catch (error) {
            addToast("error", "Failed to save settings", "Error");
            console.error(error);
        } finally {
            addToast("info", "Settings saved successfully!", "Settings Saved");
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
            window.location.href = "/login";
            return;
        }
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
            <ToggleTheme />
        </div>
    </header>

    <div class="settings-content">
        <div class="tabs">
            <button
                class={activeTab === "profile" ? "active" : ""}
                onclick={() => setActiveTab("profile")}
            >
                Profile
            </button>
            <button
                class={activeTab === "password" ? "active" : ""}
                onclick={() => setActiveTab("password")}
            >
                Password
            </button>
            <button
                class={activeTab === "deactivate" ? "active" : ""}
                onclick={() => setActiveTab("deactivate")}
            >
                Account
            </button>
            <button
                class={activeTab === "settings" ? "active" : ""}
                onclick={() => setActiveTab("settings")}
            >
                Other
            </button>
            <button
                class={activeTab === "notification" ? "active" : ""}
                onclick={() => setActiveTab("notification")}
            >
                Notifications
            </button>
            <button
                class={activeTab === "blocked" ? "active" : ""}
                onclick={() => setActiveTab("blocked")}
            >
                Blocked
            </button>
        </div>

        <div class="tab-content">
            {#if activeTab === "profile"}
                <form onsubmit={handleUpdateProfile}>
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
                <form onsubmit={handleChangePassword}>
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

                <form onsubmit={handleDeactivateAccount}>
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
                            disabled={isLoading}
                            style="background-color: darkred;"
                        >
                            {isLoading ? "Processing..." : "Deactivate Account"}
                        </button>
                    </div>
                </form>
            {:else if activeTab === "settings"}
                <form onsubmit={saveOtherSettings}>
                    <div class="form-group">
                        <label for="font-size">Font Size</label>
                        <select
                            id="font-size"
                            bind:value={otherSettings.font_size}
                        >
                            <option value="small">Small</option>
                            <option value="medium">Medium</option>
                            <option value="large">Large</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="font-color">Font Color</label>
                        <select
                            id="font-color"
                            bind:value={otherSettings.font_color}
                        >
                            <option value="default">Black</option>
                        </select>
                    </div>

                    <div class="checkbox-group">
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
            {:else if activeTab === "notification"}
                <form onsubmit={saveNotificationSettings}>
                    <div class="checkbox-group">
                        <label for="notif_like">Like Notifications</label>
                        <input
                            type="checkbox"
                            id="notif_like"
                            bind:checked={notificationSettings.notif_like}
                        />
                    </div>
                    <div class="checkbox-group">
                        <label for="notif_community"
                            >Community Notifications</label
                        >
                        <input
                            type="checkbox"
                            id="notif_community"
                            bind:checked={notificationSettings.notif_community}
                        />
                    </div>
                    <div class="checkbox-group">
                        <label for="notif_repost">Repost Notifications</label>
                        <input
                            type="checkbox"
                            id="notif_repost"
                            bind:checked={notificationSettings.notif_repost}
                        />
                    </div>
                    <div class="checkbox-group">
                        <label for="notif_follow">Follow Notifications</label>
                        <input
                            type="checkbox"
                            id="notif_follow"
                            bind:checked={notificationSettings.notif_follow}
                        />
                    </div>
                    <div class="checkbox-group">
                        <label for="notif_follow"
                            >Newsletter Notifications</label
                        >
                        <input
                            type="checkbox"
                            id="notif_newsletter"
                            bind:checked={notificationSettings.notif_newsletter}
                        />
                    </div>

                    <div class="form-actions">
                        <button type="submit" disabled={isLoading}>
                            {isLoading ? "Saving..." : "Save Settings"}
                        </button>
                    </div>
                </form>
            {:else if activeTab === "blocked"}
                <div><BlockedUser /></div>
            {/if}
        </div>
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/settings.scss";
</style>
