export function formatKopecks(
  kopecks: number,
): string {
  const rubles = kopecks / 100

  return new Intl.NumberFormat(
    'ru-RU',
    {
      style: 'currency',
      currency: 'RUB',
    },
  ).format(rubles)
}