<script lang="ts">
    const limits = [25, 30, 35];

    let currLimit = $state(25);
    let currentPage = $state(1);
    let { totalItems }: { totalItems: number } = $props();
    let totalPages = $state(0);
    let onPageChange: (page: number, limit: number) => void = () => {};

    $effect(() => {
        totalPages = Math.max(Math.ceil(totalItems / currLimit), 1);
    });
</script>

<div class="pagination">
    <div class="left">
        <span>Show: </span>
        <select
            bind:value={currLimit}
            onchange={() => {
                currentPage = 1;
                onPageChange(currentPage, currLimit);
            }}
        >
            {#each limits as l}
                <option value={l}>{l}</option>
            {/each}
        </select>
    </div>
    <div class="center">
        <button
            onclick={() =>
                currentPage > 1 && onPageChange(--currentPage, currLimit)}
            disabled={currentPage === 1}
        >
            Previous
        </button>
        <span>Page {currentPage} of {totalPages}</span>
        <button
            onclick={() =>
                currentPage < totalPages &&
                onPageChange(++currentPage, currLimit)}
            disabled={currentPage === totalPages}
        >
            Next
        </button>
    </div>
    <div class="right">
        <span>{totalItems} items total</span>
    </div>
</div>

<style lang="scss">
    .pagination {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.5rem;
        background: var(--background);
        color: var(--text);

        .left,
        .center,
        .right {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            fill: var(--primary);
        }

        select {
            background: var(--secondary);
            color: var(--text);
            border: 1px solid #ccc;
            padding: 0.25rem;
            border-radius: 4px;
        }

        button {
            background: var(--primary);
            color: white;
            border: none;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            cursor: pointer;

            &:hover:not(:disabled) {
                background: var(--primary-hover);
            }
            &:disabled {
                background: gray;
                cursor: not-allowed;
            }
        }
    }
</style>
