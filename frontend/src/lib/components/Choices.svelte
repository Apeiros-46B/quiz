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

<div class="container">
    {#each choices as choice, cn}
        <button class:active={$active == cn} on:click={onClick(cn)}>{choice}</button>
    {/each}
</div>

<style>
    .container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: space-around;
    }

    button {
        margin-bottom: 8px;

        font-size: 1rem;

        background-color: var(--bg3);

        box-shadow: 0px 0px transparent;
    }

    button {
        width: 10rem;
        max-width: 10rem;
    }

    /* {{{ adaptive button width based on screen width */
    @media (min-width: 600px) {
        button {
            width: 20rem;
            max-width: 20rem;
        }
    }

    @media (min-width: 700px) {
        button {
            width: 24rem;
            max-width: 24rem;
        }
    }

    @media (min-width: 800px) {
        button {
            width: 26rem;
            max-width: 26rem;
        }
    }

    @media (min-width: 900px) {
        button {
            width: 28rem;
            max-width: 28rem;
        }
    }

    @media (min-width: 1000px) {
        button {
            width: 30rem;
            max-width: 30rem;
        }
    }
    /* }}} */

    button.active {
        color: var(--bg1);
        background-color: var(--accent);
    }
</style>
