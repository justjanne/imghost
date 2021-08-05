import {createContext, FC, FunctionComponent, PropsWithChildren, useCallback, useContext, useState} from "react";
import {createPortal} from "react-dom";

type ChildrenComponent = FunctionComponent<PropsWithChildren<{}>>;

const ErrorContext = createContext<ChildrenComponent | null>(null);

export function useErrorDisplay(): [FC, ChildrenComponent] {
    const [ref, setRef] = useState<HTMLDivElement | null>(null);
    const errorDisplay = useCallback(
        () => (
            <div ref={setRef}/>
        ),
        []
    );
    const errorRenderer: ChildrenComponent = useCallback(
        ({children}: PropsWithChildren<{}>) => ref && createPortal(
            children,
            ref
        ),
        [ref]
    );
    const errorWrapper = useCallback(
        ({children}: PropsWithChildren<{}>) => (
            <ErrorContext.Provider value={errorRenderer}>
                {children}
            </ErrorContext.Provider>
        ),
        [errorRenderer]
    );

    return [errorDisplay, errorWrapper];
}

export const ErrorPortal: ChildrenComponent = (props: PropsWithChildren<{}>) => {
    const errorRenderer = useContext(ErrorContext);
    if (!errorRenderer) {
        return null;
    }

    return errorRenderer(props);
}
