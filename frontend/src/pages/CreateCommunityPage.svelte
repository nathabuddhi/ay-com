<script lang="ts">
    import { onMount } from "svelte";
    import {
        createNewCommunity,
        getCommunityCategories,
    } from "../controllers/community-controller";
    import { addToast } from "../stores/toast-wrapper";

    let isSubmitting = $state(false);
    let name = $state("");
    let description = $state("");
    let rules = $state("");
    let categories = $state<string[] | null>(null);
    let selectedCategories = $state<string[]>([]);
    let icon: File | null = $state(null);
    let iconPreview: string | null = $state(null);
    let banner: File | null = $state(null);
    let bannerPreview: string | null = $state(null);

    onMount(async () => {
        const response = await getCommunityCategories();

        if (response.success) {
            categories = response.payload?.categories || ["General"];
        } else {
            addToast("error", "Failed to load categories: ");
            categories = ["General"];
        }
    });

    function handleIconChange(event: Event): void {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files[0]) {
            icon = input.files[0];
            iconPreview = URL.createObjectURL(icon);
        }
    }

    function handleBannerChange(event: Event): void {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files[0]) {
            banner = input.files[0];
            bannerPreview = URL.createObjectURL(banner);
        }
    }

    async function handleSubmit(event: Event): Promise<void> {
        event.preventDefault();
        isSubmitting = true;

        if (!name || !description || !rules || !icon || !banner) {
            addToast("error", "Please fill in all fields and upload images.");
            isSubmitting = false;
            return;
        }
        if (selectedCategories.length === 0) {
            addToast("error", "Please select at least one category.");
            isSubmitting = false;
            return;
        }

        const formData = new FormData();
        formData.append("name", name);
        formData.append("description", description);
        formData.append("rules", rules);
        formData.append("icon", icon);
        formData.append("banner", banner);
        formData.append("category_count", String(selectedCategories.length));
        selectedCategories.forEach((category, idx) => {
            formData.append(`category_${idx + 1}`, category);
        });

        const response = await createNewCommunity(formData);

        if (response.success) {
            addToast("success", "Community created successfully!");
        } else {
            addToast(
                "error",
                `Failed to create community: ${response.message}`
            );
        }
        isSubmitting = false;
        name = "";
        description = "";
        rules = "";
        icon = null;
        iconPreview = null;
        banner = null;
        bannerPreview = null;
    }
</script>

<div>
    <h1 class="form-title">Create your account</h1>

    <form class="create-community-form" onsubmit={handleSubmit}>
        <div class="form-group">
            <label for="name" class="form-label">Community Name</label>
            <input
                type="name"
                id="name"
                class="form-input"
                bind:value={name}
                required
            />
        </div>

        <div class="form-group">
            <label for="description" class="form-label">Description</label>
            <textarea
                id="description"
                class="form-input"
                bind:value={description}
                required
            ></textarea>
        </div>

        <div class="form-group">
            <label class="form-label" for="category">
                Community Categories
            </label>
            <div class="checkbox-group">
                {#if categories}
                    {#each categories as category}
                        <label class="checkbox-label">
                            {category}
                            <input
                                type="checkbox"
                                value={category}
                                checked={selectedCategories.includes(category)}
                                onchange={() => {
                                    if (selectedCategories.includes(category)) {
                                        selectedCategories =
                                            selectedCategories.filter(
                                                (c) => c !== category
                                            );
                                    } else {
                                        selectedCategories = [
                                            ...selectedCategories,
                                            category,
                                        ];
                                    }
                                }}
                            />
                        </label>
                    {/each}
                {:else}
                    <span>Loading categories...</span>
                {/if}
            </div>
        </div>

        <div class="form-group">
            <label for="rules" class="form-label">Rules</label>
            <textarea id="rules" class="form-input" bind:value={rules} required
            ></textarea>
        </div>

        <div class="form-group">
            <label class="form-label" for="icon">Community Icon</label>
            <div class="file-input-container">
                <label for="icon" class="file-input-label">
                    {icon ? icon.name : "Choose a profile picture"}
                </label>
                <input
                    type="file"
                    id="icon"
                    class="file-input"
                    accept="image/*"
                    onchange={handleIconChange}
                />
                {#if iconPreview}
                    <img
                        src={iconPreview || "/placeholder.svg"}
                        alt="Icon preview"
                        class="file-preview visible"
                    />
                {/if}
            </div>
        </div>

        <div class="form-group">
            <label class="form-label" for="banner">Community Banner</label>
            <div class="file-input-container">
                <label for="banner" class="file-input-label">
                    {banner ? banner.name : "Choose a profile picture"}
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
        </div>

        <button type="submit" class="register-button" disabled={isSubmitting}>
            {isSubmitting ? "Creating community..." : "Create Community"}
        </button>
    </form>
</div>

<style lang="scss">
    @use "../styles/create-community.scss";
</style>
