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

    import { submitToBackend } from '$lib/app/submit';
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
        setDisplayByID('questionContainer', 'block');

        timeStart = Date.now();
    }
    // }}}

    // {{{ get the current question
    function getQuestion() {
        let q = data.questions[qn];
        qID = q.id;
        question = q.question;
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
        let result = await submitToBackend();
        correct = result.correct;
        keyword = result.keyword;

        if (correct) {
            // embed stats in keyword
            keyword = keyword.replace('%userID%', $userID);
            keyword = keyword.replace('%attempts%', $attemptCount.toString());
            keyword = keyword.replace('%attemptTime%', $attemptTime.toString());
            keyword = keyword.replace('%totalTime%', $totalTime.toString());

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

<!-- start screen -->
<div class="container" id="startContainer">
    {@html data.startHTML}

    <button on:click={startQuiz}>{data.startButtonTxt}</button>
</div>

<!-- question screen -->
<div class="container" id="questionContainer" style="display: none">
    <button on:click={next}>{nextTxt}</button>

    <h1 class="question">{question}</h1>
    <Choices {qID} {choices} />
</div>

<!-- end screen -->
<div class="container" id="endContainer" style="display: none">
    {@html correct ? data.successHTML : data.failureHTML}
    {@html keyword}
</div>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Work+Sans:wght@400;500;600;700&display=swap');

    :global(body) {
        /* {{{ variables */
        --gray0: #282c34;
        --gray1: #2b3339;
        --gray2: #323c41;
        --gray3: #3a454a;
        --gray4: #445055;
        --gray5: #607279;
        --gray6: #7a8487;
        --gray7: #859289;
        --gray8: #9da9a0;

        --white: #d3c6aa;

        --red: #e67e80;
        --orange: #e69875;
        --yellow: #ddbc7f;
        --green: #a7c080;
        --teal: #83c092;
        --blue: #7fbbb3;
        --purple: #d699b6;

        --visual_bg: #503946;
        --bg_yellow: #4a4940;
        --diff_del: #4e3e43;
        --diff_add: #404d44;
        --diff_mod: #394f5a;
        /* }}} */

        /* {{{ center elements */
        text-align: center;
        align-items: center;
        justify-content: center;

        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        -ms-transform: translate(-50%, -50%);
        /* }}} */

        font-family: 'Work Sans', sans-serif;
        font-size: 1.5rem;
        color: var(--white);
        padding: 0 0.5rem;
        background: var(--gray1);
        line-height: 1;
        max-width: 30rem;
        margin: auto;
    }

    :global(h1) {
        font-size: 2.2rem !important;
        margin-top: 1rem;
        margin-bottom: 1rem;
    }

    :global(strong) {
        font-size: 1.2rem !important;
    }

    /* {{{ buttons */
    :global(button) {
        font-family: inherit;
        font-weight: 500;
        font-size: 1.2rem;

        text-align: center;
        text-decoration: none;

        display: inline-block;
        padding: 4px 8px;
        border: none;
        border-radius: 8px;

        color: var(--gray7);
        background-color: var(--gray2);
        transition: 0.2s all;
        transition-property: color, background-color;
    }

    :global(button):hover {
        color: var(--white);
        background-color: var(--diff_add);
    }

    :global(button):active {
        color: var(--gray1);
        background-color: var(--green);
    }
    /* }}} */
</style>
