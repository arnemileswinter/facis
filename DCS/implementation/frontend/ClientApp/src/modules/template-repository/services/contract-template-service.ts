import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

const http = axios.create({
    baseURL: API_BASE_URL,
    headers: { 'Content-Type': 'application/json' }
})

export const ContractTemplateService = {
    create(data: any) {
        return http.post('/template/create', data).then(res => res.data)
    },

    submit(data: any) {
        return http.post('/template/submit', data).then(res => res.data)
    },

    update(data: string) {
        return http.put('/template/update', data).then(res => res.data)
    },

    retrieve() {
        return http.get('/template/retrieve')
            .then(res => Array.isArray(res.data) ? res.data : [])
            .catch(err => {
                console.error('Retrieve Error:', err)
                return []
            })
    },

    retrieveById(templateId: string) {
        return http.get(`/template/retrieve/${templateId}`).then(res => res.data)
    },

    approve(data: any) {
        return http.post('/template/approve', data).then(res => res.data)
    },

    reject(data: any) {
        return http.post('/template/reject', data).then(res => res.data)
    }
}
