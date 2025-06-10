<script lang="ts">
    import { Image, ListPlus, X } from "@lucide/svelte";
    import { AVATAR_IMG } from "../env_var";
    import { addToast } from "../stores/toast-wrapper";
    import {
        getCategories,
        postThread,
    } from "../controllers/thread-controller";
    import { onMount } from "svelte";

    const maxWords: number = 100;
    let postText: string = $state("");

    let wordCount = $derived(
        postText
            .trim()
            .split(/\s+/)
            .filter((word) => word.length > 0).length
    );

    let progress = $derived(Math.min((wordCount / maxWords) * 100, 100));

    let circleColor = $derived(
        wordCount >= maxWords
            ? "#ff0000"
            : wordCount > maxWords * 0.8
              ? "#ffa500"
              : "#00ff00"
    );

    let selectedFiles: FileList | null = $state<FileList | null>(null);
    let previewImages: string[] = $state<string[]>([]);
    let pollOptions: string[] = $state([]);
    let showPoll = $state(false);

    let categories: string[] = $state<string[]>(["LOADING"]);
    let selectedCategory: string = $state("");
    const permissions: string[] = [
        "Everyone",
        "Accounts You Follow",
        "Verified Accounts",
    ];
    let selectedPermission: string = $state("");
    let isScheduled: boolean = $state(false);
    let scheduledDate = $state("");
    let scheduledTime = $state("");

    let scheduledAt = $derived(
        isScheduled && scheduledDate && scheduledTime
            ? `${scheduledDate} ${scheduledTime}:00`
            : ""
    );

    let {
        isOpen = $bindable(),
        mode = "post",
        threadId = $bindable(),
    }: {
        isOpen: boolean;
        mode: "post" | "comment";
        threadId?: string;
    } = $props();

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
            console.log(scheduledAt);
            if (mode === "post") {
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
            }
            let response;
            if (mode === "comment" && threadId) {
                response = await postThread(
                    postText,
                    "Reply",
                    "Public",
                    selectedFiles,
                    [],
                    threadId,
                    ""
                );
            } else if (mode === "post") {
                response = await postThread(
                    postText,
                    selectedCategory,
                    selectedPermission,
                    selectedFiles,
                    pollOptions,
                    "",
                    scheduledAt
                );
            } else {
                throw new Error("Invalid mode specified.");
            }

            if (!response.success) {
                throw new Error(response.message);
            }

            addToast(
                "success",
                `${mode === "comment" ? "Comment" : "Thread"} created successfully!`,
                "Success!"
            );

            isOpen = false;
            postText = "";
            selectedFiles = null;
            previewImages = [];
            pollOptions = [];
            showPoll = false;
            selectedCategory = "";
            selectedPermission = "";
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } catch (error) {
            addToast(
                "error",
                `Failed to ${mode === "post" ? "create post" : "post comment"}: ` +
                    (error as Error).message,
                "Error!"
            );
        }
    }

    onMount(async () => {
        const response = await getCategories();
        if (response.success && response.payload) {
            categories = response.payload;
        } else {
            categories = [
                "General",
                "Technology",
                "Health",
                "Entertainment",
                "Other",
                "Entertainment",
                "Education",
            ];
        }
    });
</script>

<div class="create-post-content">
    <div class="create-post-avatar">
        <img
            src={`${AVATAR_IMG}/${localStorage.getItem("user_id") || "1"}.png`}
            alt="Your avatar"
        />
    </div>

    <textarea
        placeholder={mode === "post"
            ? "What's happening?"
            : "Write a comment..."}
        bind:value={postText}
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
                    <button class="remove-image" onclick={() => removeImage(i)}>
                        <X />
                    </button>
                </div>
            {/each}
        </div>
    {/if}

    {#if showPoll && mode === "post"}
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
                        <button onclick={() => removePollOption(index)}
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
        <div>
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
                {#if mode === "post"}
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
                        <label for="permission-select">Permissions</label>
                        <select
                            id="permission-select"
                            bind:value={selectedPermission}
                        >
                            <option value="" disabled selected
                                >Select a permission</option
                            >
                            {#each permissions as permission}
                                <option value={permission}>{permission}</option>
                            {/each}
                        </select>
                    </div>
                    <div class="counter-container">
                        <svg width="50" height="50" viewBox="0 0 100 100">
                            <circle
                                cx="50"
                                cy="50"
                                r="45"
                                fill="none"
                                stroke="#e0e0e0"
                                stroke-width="10"
                            />
                            <circle
                                cx="50"
                                cy="50"
                                r="45"
                                fill="none"
                                stroke={circleColor}
                                stroke-width="10"
                                stroke-dasharray={`${progress * 2.83} ${283 - progress * 2.83}`}
                                transform="rotate(-90 50 50)"
                            />
                            <text
                                x="50"
                                y="55"
                                text-anchor="middle"
                                font-size="20"
                                fill="#333"
                                class="word-count"
                            >
                                {wordCount}/{maxWords}
                            </text>
                        </svg>
                        {#if wordCount > maxWords}
                            <p class="warning">
                                You have exceeded the word limit!
                            </p>
                        {/if}
                    </div>
                {/if}
            </div>
            {#if mode === "post"}
                <div class="schedule-selector">
                    <label for="schedule-select">Schedule</label>
                    <input
                        type="checkbox"
                        id="schedule-checkbox"
                        bind:checked={isScheduled}
                    />
                    {#if isScheduled}
                        <input type="date" bind:value={scheduledDate} />
                        <input type="time" bind:value={scheduledTime} />
                    {/if}
                </div>
            {/if}
        </div>

        <button
            class="post-button"
            disabled={(!postText &&
                previewImages.length === 0 &&
                (!showPoll || pollOptions.every((p) => !p.trim()))) ||
                wordCount > maxWords}
            onclick={handleSubmit}
        >
            {mode === "post" ? "Post" : "Comment"}
        </button>
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/create-post.scss";
</style>
