<script lang="ts">
    // stores
    import {
        answers,
        active,
        attemptCount,
        attemptTime,
        totalTime,
        userID,
    } from '$lib/app/stores';

    import Choices from '$lib/components/Choices.svelte';

    // load settings & questions
    import type { PageData } from './$types';
    export let data: PageData;

    // helpers
    function setDisplayByID(id: string, display: string) {
        const element = document.getElementById(id);
        if (element != null) element.style.display = display;
    }

    function preventNavigation(event: Event) {
        event.preventDefault()
        event.returnValue = false
    }

    // {{{ state
    let timeStart = 0;

    let inQuiz = false;
    let qID = '';
    let qn = 0;
    let question = '';
    let choices: string[] = [];

    let nextTxt: string = data.nextButtonTxt;
    let correct: boolean = false;
    let keyword: string = '';
    // }}}

    // {{{ start the quiz
    function startQuiz() {
        // prevent navigating away
        window.addEventListener('beforeunload', preventNavigation);
        inQuiz = true;

        getQuestion();

        setDisplayByID('startContainer', 'none');
        setDisplayByID('questionContainer', 'flex');

        timeStart = Date.now();
    }
    // }}}

    // {{{ get the current question
    function getQuestion() {
        let q = data.questions[qn];
        qID = q.id;
        question = q.question;
        console.log(question);
        choices = q.choices;
    }
    // }}}

    // {{{ go to next question
    function next() {
        // confirm if not answered
        if ($answers[qID] == undefined) {
            if (!confirm(data.proceedConfirmTxt)) return;
        }

        if (qn < data.questions.length - 1) {
            // deselect button
            $active = -1;

            qn++;
            getQuestion();
            if (qn == data.questions.length - 1) nextTxt = data.submitButtonTxt;
        } else {
            submit();
        }
    }
    // }}}

    // {{{ submit quiz
    async function submit() {
        // update stats
        $attemptTime = Date.now() - timeStart;
        $totalTime += $attemptTime;
        $attemptCount++;

        // allow navigating away
        window.removeEventListener('beforeunload', preventNavigation);
        inQuiz = false;

        // check answers
        let attemptTimeStr = $attemptTime.toString();
        let totalTimeStr = $totalTime.toString();
        let result = await fetch('/api/submit', {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
                'X-User-ID': $userID,
                'X-Attempt-Time': attemptTimeStr,
                'X-Total-Time': totalTimeStr,
            },
            body: JSON.stringify($answers),
        }).then(response => response.json());

        correct = result.correct;
        keyword = result.keyword;

        if (correct) {
            // embed stats in keyword
            keyword = keyword.replace('%userID%', $userID);
            keyword = keyword.replace('%attempts%', $attemptCount.toString());
            keyword = keyword.replace('%attemptTime%', attemptTimeStr);
            keyword = keyword.replace('%totalTime%', totalTimeStr);

            // reset stats
            $attemptCount = 0;
            $totalTime = 0;
        }

        setDisplayByID('questionContainer', 'none');
        setDisplayByID('endContainer', 'block');
    }
    // }}}
</script>

<svelte:head>
    <title>{inQuiz ? `Quiz - ${question}` : 'Quiz'}</title>
</svelte:head>

<div class="rootContainer">
    <!-- start screen -->
    <div class="container" id="startContainer" style="display: flex">
        <div class="card">
            {@html data.startHTML}
        </div>

        <br>
        <button on:click={startQuiz}>{data.startButtonTxt}</button>
    </div>

    <!-- question screen -->
    <div class="container" id="questionContainer" style="display: none">
        <div class="card">
            <p class="question"><strong>{question}</strong></p>
            <Choices {qID} {choices} />
        </div>

        <br>
        <button on:click={next}>{nextTxt}</button>
    </div>

    <!-- end screen -->
    <div class="container" id="endContainer" style="display: none">
        <div class="card">
            {@html correct ? data.successHTML : data.failureHTML}
            {@html keyword}
        </div>
    </div>
</div>

<style>
    /* @import url('https://fonts.googleapis.com/css2?family=Work+Sans:wght@400;500;600;700&display=swap'); */

    :global(body) {
        /* {{{ variables */
        --fg:       #d3c6aa;
        --fg-muted: #859289;

        --accent:    #a7c080;
        --bg-accent: #404d44;

        --bg1: #2b3339;
        --bg2: #323c41;
        --bg3: #3a454a;
        --bg4: #445055;
        /* }}} */

        /* {{{ center elements */
        /* text-align: center; */
        /* align-items: center; */
        /* justify-content: center; */

        /* position: relative; */
        /* top: 50%; */
        /* left: 50%; */
        /* transform: translateY(50%); */
        /* transform: translate(0, 25%); */
        /* -ms-transform: translate(-50%, -50%); */
        /* }}} */

        font-family: 'Work Sans', sans-serif;
        font-size: 1rem;
        line-height: 1.2;

        color: var(--fg);
        background-color: var(--bg1);

        margin: auto;
        margin-top: 8vh;
    }

    /* {{{ markup elements */
    :global(h1) {
        font-size: 1.6rem;
        margin-top: 0px;
        margin-bottom: 1rem;
    }

    :global(strong) {
        font-size: 1.2rem;
        font-weight: 700;
    }

    .question {
        font-size: 1.2rem;
        margin-block-start: 0px;
        margin-block-end: 1em;
    }
    /* }}} */

    /* {{{ containers */
    .rootContainer {
        display: grid;
        height: 100%;
        width: 100%;
    }

    .container {
        display: flex;
        flex-direction: column;
        box-sizing: border-box;
        max-height: 84vh;
        margin: auto;

        text-align: center;
        align-items: center;
    }

    /* {{{ adaptive container width based on screen width */
    .container {
        max-width: 72vw;
    }

    @media (min-width: 600px) {
        .container {
            max-width: 68vw;
        }
    }

    @media (min-width: 700px) {
        .container {
            max-width: 64vw;
        }
    }

    @media (min-width: 800px) {
        .container {
            max-width: 60vw;
        }
    }

    @media (min-width: 900px) {
        .container {
            max-width: 56vw;
        }
    }

    @media (min-width: 1000px) {
        .container {
            max-width: 48vw;
        }
    }
    /* }}} */

    .card {
        flex-grow: 1;
        padding: 20px;
        overflow-wrap: break-word;
        overflow-y: auto;

        text-align: left;

        background-color: var(--bg2);

        border-radius: 8px;
        box-shadow: 0 0 8px var(--bg1);
    }
    /* }}} */

    /* {{{ scrollbar */
    ::-webkit-scrollbar {
        width: 16px;
    }

    ::-webkit-scrollbar-track {
        opacity: 0px;
    }

    ::-webkit-scrollbar-thumb {
        background-color: var(--bg3);
        background-clip: content-box;

        border: 6px solid transparent;
        border-radius: 16px;
    }

    ::-webkit-scrollbar-thumb:hover {
        background-color: var(--bg4);
    }
    /* }}} */

    /* {{{ buttons */
    :global(button) {
        display: inline-block;
        max-width: 10em;
        padding: 4px 8px;

        text-align: center;
        text-decoration: none;

        font-family: inherit;
        font-weight: 500;
        font-size: 1.2rem;

        color: var(--fg-muted);
        background-color: var(--bg2);

        border: none;
        border-radius: 8px;
        box-shadow: 0 0 8px var(--bg1);

        transition: 0.2s all;
        transition-property: color, background-color;
    }

    :global(button):hover {
        color: var(--fg);
        background-color: var(--bg-accent);
    }

    :global(button):active {
        color: var(--bg1);
        background-color: var(--accent);
    }
    /* }}} */
</style>
