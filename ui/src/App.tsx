import React from 'react';
import './App.css';
import ImageList from "./components/ImageList";
import {BaseUrlProvider} from './api/baseUrlContext';
import {QueryClient, QueryClientProvider} from "react-query";

const queryClient = new QueryClient();

function App() {
    return (
        <BaseUrlProvider value="http://localhost:8080/">
            <QueryClientProvider client={queryClient}>
                <ImageList/>
            </QueryClientProvider>
        </BaseUrlProvider>
    );
}

export default App;
