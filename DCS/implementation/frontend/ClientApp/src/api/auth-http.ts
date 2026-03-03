import axios from 'axios'

const API_BASE_URL = import.meta.env.DCS_API_URL

const http = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
})

export default http
