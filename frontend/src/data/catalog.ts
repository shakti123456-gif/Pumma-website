import type { Product } from './products'

export const SPORTS = ['Running', 'Training', 'Football', 'Basketball', 'Lifestyle'] as const
export const GENDERS = ['Men', 'Women', 'Unisex'] as const
export const PRODUCT_TYPES = ['Clothing', 'Shoes'] as const

export type Sport = (typeof SPORTS)[number]
export type Gender = (typeof GENDERS)[number]
export type ProductType = (typeof PRODUCT_TYPES)[number]

export function filterBySlug(products: Product[], slug?: string): Product[] {
  if (!slug || slug === 'all' || slug === 'new') return products
  if (slug === 'sale') return products.filter(p => p.originalPrice)

  const key = slug.toLowerCase()
  const sport = SPORTS.find(s => s.toLowerCase() === key)
  if (sport) return products.filter(p => p.sport === sport)

  if (key === 'men') return products.filter(p => p.gender === 'Men')
  if (key === 'women') return products.filter(p => p.gender === 'Women')
  if (key === 'clothing') return products.filter(p => p.type === 'Clothing')
  if (key === 'shoes') return products.filter(p => p.type === 'Shoes')

  return products
}

export function slugTitle(slug?: string): string {
  if (!slug || slug === 'all') return 'Shop all'
  if (slug === 'new') return 'New arrivals'
  if (slug === 'sale') return 'Sale'

  const sport = SPORTS.find(s => s.toLowerCase() === slug.toLowerCase())
  if (sport) return sport

  if (slug.toLowerCase() === 'men') return "Men's"
  if (slug.toLowerCase() === 'women') return "Women's"
  if (slug.toLowerCase() === 'clothing') return 'Clothing'
  if (slug.toLowerCase() === 'shoes') return 'Shoes'

  return slug.replace(/^./, c => c.toUpperCase())
}

export function applyFilters(
  products: Product[],
  filters: { sports: Sport[]; genders: Gender[]; types: ProductType[] }
): Product[] {
  return products.filter(p => {
    if (filters.sports.length && !filters.sports.includes(p.sport)) return false
    if (filters.genders.length && !filters.genders.includes(p.gender)) return false
    if (filters.types.length && !filters.types.includes(p.type)) return false
    return true
  })
}
