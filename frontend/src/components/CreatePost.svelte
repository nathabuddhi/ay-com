<script lang="ts">
    let postText = "";
    let selectedFiles: FileList | null = null;
    let previewImages: string[] = [];

    function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files) {
            selectedFiles = input.files;
            previewImages = [];

            for (let i = 0; i < selectedFiles.length; i++) {
                const file = selectedFiles[i];
                const reader = new FileReader();

                reader.onload = (e) => {
                    if (e.target?.result) {
                        previewImages = [
                            ...previewImages,
                            e.target.result as string,
                        ];
                    }
                };

                reader.readAsDataURL(file);
            }
        }
    }

    function removeImage(index: number) {
        previewImages = previewImages.filter((_, i) => i !== index);

        if (selectedFiles) {
            const dt = new DataTransfer();
            for (let i = 0; i < selectedFiles.length; i++) {
                if (i !== index) {
                    dt.items.add(selectedFiles[i]);
                }
            }
            selectedFiles = dt.files;
        }
    }

    function handleSubmit() {
        console.log("Post text:", postText);
        console.log("Selected files:", selectedFiles);

        postText = "";
        selectedFiles = null;
        previewImages = [];
    }
</script>

<div class="create-post">
    <div class="create-post-avatar">
        <img
            src={`http://localhost:5000/images/profile/${localStorage.getItem("user_id") || "1"}`}
            alt="Your avatar"
        />
    </div>

    <div class="create-post-content">
        <textarea placeholder="What's happening?" bind:value={postText} rows="3"
        ></textarea>

        {#if previewImages.length > 0}
            <div
                class="image-previews {previewImages.length > 1
                    ? 'multiple-images'
                    : ''}"
            >
                {#each previewImages as preview, i}
                    <div class="preview-container">
                        <img
                            src={preview || "/placeholder.svg"}
                            alt="Preview"
                        />
                        <button
                            class="remove-image"
                            on:click={() => removeImage(i)}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <line x1="18" y1="6" x2="6" y2="18"></line>
                                <line x1="6" y1="6" x2="18" y2="18"></line>
                            </svg>
                        </button>
                    </div>
                {/each}
            </div>
        {/if}

        <div class="create-post-actions">
            <div class="media-actions">
                <label class="media-button">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"
                        ></rect>
                        <circle cx="8.5" cy="8.5" r="1.5"></circle>
                        <polyline points="21 15 16 10 5 21"></polyline>
                    </svg>
                    <input
                        type="file"
                        accept="image/*"
                        multiple
                        on:change={handleFileChange}
                    />
                </label>

                <button class="media-button">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <path
                            d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                        ></path>
                        <polyline points="14 2 14 8 20 8"></polyline>
                        <line x1="16" y1="13" x2="8" y2="13"></line>
                        <line x1="16" y1="17" x2="8" y2="17"></line>
                        <polyline points="10 9 9 9 8 9"></polyline>
                    </svg>
                </button>

                <button class="media-button">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <circle cx="12" cy="12" r="10"></circle>
                        <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
                        <line x1="9" y1="9" x2="9.01" y2="9"></line>
                        <line x1="15" y1="9" x2="15.01" y2="9"></line>
                    </svg>
                </button>

                <button class="media-button">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"
                        ></path>
                    </svg>
                </button>
            </div>

            <button
                class="post-button"
                disabled={!postText && previewImages.length === 0}
                on:click={handleSubmit}
            >
                Post
            </button>
        </div>
    </div>
</div>

<style lang="scss">
    @use "../styles/home.scss";
</style>
