<script lang="ts">
    import { answers, active } from '$lib/app/stores';

    export let qID: string;
    export let choices: string[];

    function onClick(cn: number) {
        return () => {
            $active = $answers[qID] = cn;
        };
    }
</script>

{#each choices as choice, cn}
    <button class:active={$active == cn} on:click={onClick(cn)}>{choice}</button>
{/each}

<style>
    button {
        font-size: 1rem;
        margin: 4px;
    }

    button.active {
        color: var(--gray1);
        background-color: var(--green);
    }
</style>
