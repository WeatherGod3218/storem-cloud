import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { BrowserRouter, Routes, Route } from "react-router";
import { LandingPage } from './pages/LandingPage.tsx';
import { ThemeProvider } from './components/ThemeProvider.tsx';
import { VideoPage } from './pages/VideoPage.tsx';
import { TestPage } from './pages/TestPage.tsx';

import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'

const queryProvider = new QueryClient()

createRoot(document.getElementById('root')!).render(
	<StrictMode>					
		<QueryClientProvider client={queryProvider}>
			<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">	
				<BrowserRouter>
					<Routes>
							<Route path="" element={<LandingPage/>}/>
							<Route path="/" element={<LandingPage/>}/>
							<Route path="/video/:id" element={<VideoPage/>}/>
							<Route path="/test" element={<TestPage/>}/>					
					</Routes>
				</BrowserRouter>
			</ThemeProvider>
		</QueryClientProvider>
  </StrictMode>
)
