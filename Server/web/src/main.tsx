
import './index.css'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import { BrowserRouter, Routes, Route } from "react-router";
import { ThemeProvider } from './components/ThemeProvider.tsx';
import { AuthProvider } from './context/AuthContext.tsx';

import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'

import { LandingPage } from './pages/LandingPage.tsx';
import { VideoPage } from './pages/VideoPage.tsx';
import { TestPage } from './pages/TestPage.tsx';
import LoginPage from './pages/LoginPage.tsx';


const queryProvider = new QueryClient()

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<AuthProvider>
			<QueryClientProvider client={queryProvider}>
				<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">	
					<BrowserRouter>
						<Routes>
								<Route path="" element={<LandingPage/>}/>
								<Route path="/" element={<LandingPage/>}/>
								<Route path="/video/:id" element={<VideoPage/>}/>
								<Route path="/test" element={<TestPage/>}/>
								<Route path="/login" element={<LoginPage/>}/>						
						</Routes>
					</BrowserRouter>
				</ThemeProvider>
			</QueryClientProvider>
		</AuthProvider>	
  </StrictMode>
)
