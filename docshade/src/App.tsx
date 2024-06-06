import React from 'react';
import Home from './pages/Home/Home';
import Header from './components/Header/Header';
import Footer from './components/Footer/Footer';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './App.css';

const App: React.FC = () => (
  <div className="App">
    <Header />
    <main className="MainContent">
      <Home />
    </main>
    <Footer />
    <ToastContainer />
  </div>
);

export default App;
