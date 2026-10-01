import axios from 'axios'

import type {
  ApiErrorResponse,
} from '../types/banking'

export function getApiErrorMessage(
  error: unknown,
): string {
  if (
    axios.isAxiosError<ApiErrorResponse>(
      error,
    )
  ) {
    return (
      error.response?.data.error.message ??
      'Не удалось выполнить запрос'
    )
  }

  return 'Произошла неизвестная ошибка'
}