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

export function parseRublesToKopecks(
  value: string,
): number | null {
  const normalized =
    value.trim().replace(',', '.')

  if (
    !/^\d+(?:\.\d{1,2})?$/.test(
      normalized,
    )
  ) {
    return null
  }

  const [
    rublesPart,
    kopecksPart = '',
  ] = normalized.split('.')

  const rubles = Number(rublesPart)

  const kopecks = Number(
    kopecksPart.padEnd(2, '0'),
  )

  if (
    !Number.isSafeInteger(rubles) ||
    !Number.isSafeInteger(kopecks)
  ) {
    return null
  }

  const total =
    rubles * 100 + kopecks

  if (
    !Number.isSafeInteger(total) ||
    total <= 0
  ) {
    return null
  }

  return total
}