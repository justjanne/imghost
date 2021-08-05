import './App.css';
import {BaseUrlProvider} from './api/baseUrlContext';
import {QueryClient, QueryClientProvider} from "react-query";
import UploadView from "./components/UploadView";
import {BrowserRouter, Link, Route} from "react-router-dom";
import ImageListPage from "./pages/ImageListPage";
import {Redirect, Switch} from "react-router";
import ImageDetailPage from "./pages/ImageDetailPage";
import {AppBar, Button, Toolbar, Typography} from "@material-ui/core";
import {useErrorDisplay} from "./components/ErrorContext";

const queryClient = new QueryClient();

function App() {
    const [ErrorDisplay, ErrorWrapper] = useErrorDisplay();

    return (
        <BaseUrlProvider value="http://localhost:8080/">
            <QueryClientProvider client={queryClient}>
                <BrowserRouter>
                    <AppBar position="static">
                        <Toolbar>
                            <Typography variant="h6">
                                i.k8r.eu
                            </Typography>
                            <Button
                                color="inherit"
                                component={Link}
                                to="/upload"
                            >
                                Upload
                            </Button>
                            <Button
                                color="inherit"
                                component={Link}
                                to="/images"
                            >
                                Images
                            </Button>
                        </Toolbar>
                    </AppBar>
                    <ErrorDisplay/>
                    <ErrorWrapper>
                        <Switch>
                            <Route path="/upload">
                                <UploadView/>
                            </Route>
                            <Route path="/images">
                                <ImageListPage/>
                            </Route>
                            <Route path="/i/:imageId">
                                <ImageDetailPage/>
                            </Route>
                            <Redirect from="/" to="/images"/>
                        </Switch>
                    </ErrorWrapper>
                </BrowserRouter>
            </QueryClientProvider>
        </BaseUrlProvider>
    );
}

export default App;
