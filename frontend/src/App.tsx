import {
  Route,
  Routes,
} from 'react-router-dom'

import {
  DashboardPage,
} from './pages/DashboardPage'

import {
  AccountPage,
} from './pages/AccountPage'

import './App.css'

function App() {
  return (
    <Routes>
      <Route
        path="/"
        element={
          <DashboardPage />
        }
      />

      <Route
        path="/accounts/:id"
        element={
          <AccountPage />
        }
      />

      <Route
        path="*"
        element={
          <main className="app">
            <section className="dashboard">
              <div className="state-message">
                Страница не найдена
              </div>
            </section>
          </main>
        }
      />
    </Routes>
  )
}

export default App