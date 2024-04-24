<script lang="ts">
	import { onMount } from "svelte";

	type SseData<T> = {
		Id: number;
		Path: string;
		Data: T;
	};

	let count: number | undefined;
	let sse: EventSource | undefined = undefined;
	onMount(() => {
		start();
	});

	const start = () => {
		sse = new EventSource("http://localhost:8080/sse");

		sse.onerror = function (e) {
            console.log(e);
            
			if (sse?.readyState === EventSource.CLOSED) {
				console.log("Connection closed by the server");
			} else {
				console.log("Error occurred");
			}
		};

		sse.addEventListener("message", (e) => {
			try {
				const res = JSON.parse(e.data) as SseData<number>;
				count = res.Data;
			} catch (error) {
				console.error(error);
				count = undefined;
			}
		});
		sse.addEventListener("open", (e) => {
			console.log("open", e);
		});
	};

	const stop = () => {
		sse?.close();
        console.log(sse);
		sse = undefined;
		count = undefined;
	};

	const inc = async () => {
		fetch("http://localhost:8080/inc", {
			method: "post",
		});
	};
</script>

<div>
	<div>Count is {count ?? ""}</div>
	<br />
	<button on:click={inc}>Inc</button>
	<br />
	<br />
	{#if !sse}
		<button on:click={start}> Start </button>
	{:else}
		<button on:click={stop}> Stop </button>
	{/if}
</div>
