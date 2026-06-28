import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('higolang_token') || '')

  const isLoggedIn = computed(() => !!token.value)

  const login = async (username, password) => {
    const response = await api.post('/api/v1/admin/login', { username, password })
    if (response.data.code === 0) {
      token.value = response.data.data.token
      localStorage.setItem('higolang_token', token.value)
      return true
    }
    throw new Error(response.data.message || '登录失败')
  }

  const logout = () => {
    token.value = ''
    localStorage.removeItem('higolang_token')
  }

  return {
    token,
    isLoggedIn,
    login,
    logout
  }
})
