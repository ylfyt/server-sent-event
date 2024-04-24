export type SseState = "connecting" | "connected" | "closed";

export class Sse<T = any> {
    private es: EventSource;
    private prevState: number;
    private stateChangeCb: ((state: SseState) => void) | undefined;
    private errorCb: ((e: Event) => void) | undefined;
    private dataCb: ((data?: T) => void) | undefined;

    constructor(url: string) {
        const eventSource = new EventSource(url);
        this.es = eventSource;
        this.prevState = eventSource.readyState;

        eventSource.addEventListener('error', this.internalOnError);
        eventSource.addEventListener('open', this.internalOnOpen);
        eventSource.addEventListener('message', this.internalOnMessage);
    }
    private internalOnError = (e: Event) => {
        this.stateTrigger();
        this.errorCb && this.errorCb(e);
    };
    private internalOnOpen = (e: Event) => {
        this.stateTrigger();
    };
    private internalOnMessage = (e: MessageEvent<any>) => {
        this.stateTrigger();

        try {
            const data = JSON.parse(e.data) as T;
            this.dataCb && this.dataCb(data);
        } catch (error) {
            console.error(error);
            this.dataCb && this.dataCb();
        }
    };

    private getState(state: number): SseState {
        return state === 0 ? 'connecting' : state === 1 ? "connected" : "closed";
    }
    onStateChange(cb: (state: SseState) => void) {
        this.stateChangeCb = cb;
        cb(this.getState(this.es.readyState));
    }
    onError(cb: (e: Event) => void) {
        this.errorCb = cb;
    }
    onData(cb: (data?: T) => void) {
        this.dataCb = cb;
    }

    close() {
        this.es.removeEventListener('error', this.internalOnError);
        this.es.removeEventListener('open', this.internalOnOpen);
        this.es.removeEventListener('message', this.internalOnMessage);
        this.es.close();
        this.stateTrigger();
    }

    private stateTrigger() {
        const curr = this.es.readyState;
        if (curr !== this.prevState) {
            this.stateChangeCb && this.stateChangeCb(this.getState(curr));
        }
        this.prevState = curr;
    }
}