import {createContext, useContext} from "react";

const BaseUrlContext = createContext<string>(window.location.href);
export const BaseUrlProvider = BaseUrlContext.Provider;
export const useBaseUrl = () => useContext<string>(BaseUrlContext);
