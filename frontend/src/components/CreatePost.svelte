<script lang="ts">
    import { Image, X, ListPlus } from "@lucide/svelte";
    import { AVATAR_IMG } from "../env_var";
    import { addToast } from "../stores/toast-wrapper";
    import { postThread } from "../controllers/thread-controller";

    let postText: string = $state("");
    let selectedFiles: FileList | null = $state<FileList | null>(null);
    let previewImages: string[] = $state<string[]>([]);
    let pollOptions: string[] = $state([]);
    let showPoll = $state(false);

    const categories = [
        "General",
        "News",
        "Gaming",
        "Education",
        "Technology",
        "Art",
    ];
    let selectedCategory = $state("");
    const permissions = ["Public", "Followers", "Friends"];
    let selectedPermission = $state("");

    let { isOpen = $bindable() }: { isOpen: boolean } = $props();

    function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (!input.files) return;

        const newFiles = Array.from(input.files);

        const dt = new DataTransfer();
        if (selectedFiles) {
            for (let i = 0; i < selectedFiles.length; i++) {
                dt.items.add(selectedFiles[i]);
            }
        }
        newFiles.forEach((file) => dt.items.add(file));
        selectedFiles = dt.files;

        const newPreviews: Promise<string>[] = newFiles.map(
            (file) =>
                new Promise((resolve) => {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        resolve(e.target?.result as string);
                    };
                    reader.readAsDataURL(file);
                })
        );

        Promise.all(newPreviews).then((results) => {
            previewImages = [...previewImages, ...results];
        });
    }

    function removeImage(index: number) {
        previewImages = previewImages.filter((_, i) => i !== index);
        if (selectedFiles) {
            const dt = new DataTransfer();
            for (let i = 0; i < selectedFiles.length; i++) {
                if (i !== index) dt.items.add(selectedFiles[i]);
            }
            selectedFiles = dt.files;
        }
    }

    function addPollOption() {
        if (pollOptions.length < 4) {
            pollOptions = [...pollOptions, ""];
        }
    }

    function removePollOption(index: number) {
        pollOptions = pollOptions.filter((_, i) => i !== index);
    }

    async function handleSubmit() {
        try {
            if (selectedCategory === "") {
                addToast("error", "Please select a category.", "Error!");
                return;
            }

            if (selectedPermission === "") {
                addToast(
                    "error",
                    "Please select a permission level.",
                    "Error!"
                );
                return;
            }

            const response = await postThread(
                postText,
                selectedCategory,
                selectedPermission,
                selectedFiles,
                pollOptions
            );

            if (!response.success) {
                throw new Error(response.message);
            }

            isOpen = false;
            addToast("success", "Thread created successfully!", "Success!");
        } catch (error) {
            addToast(
                "error",
                "Failed to create post: " + (error as Error).message,
                "Error!"
            );
        } finally {
        }
    }
</script>

{#if isOpen}
    <div class="create-post-background">
        <div class="create-post">
            <button class="close-button" onclick={() => (isOpen = false)}>
                <X />
            </button>

            <div class="create-post-avatar">
                <img
                    src={`${AVATAR_IMG}/${localStorage.getItem("user_id") || "1"}.png`}
                    alt="Your avatar"
                />
            </div>

            <div class="create-post-content">
                <textarea placeholder="What's happening?" bind:value={postText}
                ></textarea>

                {#if previewImages.length > 0}
                    <div
                        class="image-previews {previewImages.length > 1
                            ? 'multiple-images'
                            : ''}"
                    >
                        {#each previewImages as preview, i}
                            <div class="preview-container">
                                <img src={preview} alt="Preview" />
                                <button
                                    class="remove-image"
                                    onclick={() => removeImage(i)}
                                >
                                    <X />
                                </button>
                            </div>
                        {/each}
                    </div>
                {/if}

                {#if showPoll}
                    <div class="poll-section">
                        <h4>Create a Poll</h4>
                        {#each pollOptions as option, index}
                            <div class="poll-option">
                                <input
                                    type="text"
                                    bind:value={pollOptions[index]}
                                    placeholder={`Option ${index + 1}`}
                                />
                                {#if pollOptions.length > 0}
                                    <button
                                        onclick={() => removePollOption(index)}
                                        >✕</button
                                    >
                                {/if}
                            </div>
                        {/each}
                        {#if pollOptions.length < 4}
                            <button class="add-option" onclick={addPollOption}
                                >Add Option</button
                            >
                        {/if}
                    </div>
                {/if}

                <div class="create-post-actions">
                    <div class="media-actions">
                        <label class="media-button">
                            <Image />
                            <input
                                type="file"
                                accept="image/*"
                                multiple
                                onchange={handleFileChange}
                            />
                        </label>
                        <button
                            class="media-button"
                            onclick={() => {
                                showPoll = !showPoll;
                            }}
                        >
                            <ListPlus />
                        </button>
                        <div class="category-selector">
                            <label for="category-select">Category</label>
                            <select
                                id="category-select"
                                bind:value={selectedCategory}
                            >
                                <option value="" disabled selected
                                    >Select a category</option
                                >
                                {#each categories as category}
                                    <option value={category}>{category}</option>
                                {/each}
                            </select>
                        </div>
                        <div class="category-selector">
                            <label for="category-select">Permissions</label>
                            <select
                                id="category-select"
                                bind:value={selectedPermission}
                            >
                                <option value="" disabled selected
                                    >Select a permission</option
                                >
                                {#each permissions as permission}
                                    <option value={permission}
                                        >{permission}</option
                                    >
                                {/each}
                            </select>
                        </div>
                    </div>

                    <button
                        class="post-button"
                        disabled={!postText &&
                            previewImages.length === 0 &&
                            (!showPoll || pollOptions.every((p) => !p.trim()))}
                        onclick={handleSubmit}
                    >
                        Post
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
