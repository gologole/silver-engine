'use client'

import { useState, useEffect } from 'react'
import { useApi } from '@/hooks/useApi'
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'

interface Statistics {
  TotalBooks: number
  TotalUsers: number
  TotalLoans: number
  PopularCategories: { Name: string; Count: number }[]
  ActiveUsers: { Name: string; LoansCount: number }[]
}

export default function Statistics() {
  const [stats, setStats] = useState<Statistics | null>(null)
  const { fetchApi, loading, error } = useApi()

  useEffect(() => {
    const loadStatistics = async () => {
      const data = await fetchApi('/statistics')
      if (data) {
        console.log('Received statistics:', data) // Логирование для дебага
        setStats(data) // Обновляем состояние с данными статистики
      }
    }
    loadStatistics() // Вызов функции загрузки статистики
  }, [])

  if (loading) return <p>Загрузка...</p> // Показываем индикатор загрузки
  if (error) return <p className="text-red-500">{error}</p> // Показываем ошибку, если она произошла
  if (!stats) return null // Если данные еще не загружены, ничего не показываем

  return (
      <div className="space-y-4">
        <h1 className="text-2xl font-bold">Статистика библиотеки</h1>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {/* Карточка с количеством книг */}
          <Card>
            <CardHeader>
              <CardTitle>Всего книг</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{stats.TotalBooks}</p>
            </CardContent>
          </Card>

          {/* Карточка с количеством пользователей */}
          <Card>
            <CardHeader>
              <CardTitle>Всего пользователей</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{stats.TotalUsers}</p>
            </CardContent>
          </Card>

          {/* Карточка с количеством выдач */}
          <Card>
            <CardHeader>
              <CardTitle>Всего выдач</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{stats.TotalLoans}</p>
            </CardContent>
          </Card>
        </div>

        {/* График популярных категорий */}
        <Card>
          <CardHeader>
            <CardTitle>Популярные категории</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={stats.PopularCategories}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="Name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="Count" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* График активных пользователей */}
        <Card>
          <CardHeader>
            <CardTitle>Активные пользователи</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={stats.ActiveUsers}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="Name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="LoansCount" fill="#82ca9d" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>
  )
}
