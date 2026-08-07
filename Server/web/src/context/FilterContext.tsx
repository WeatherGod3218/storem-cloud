import { createContext, useContext, useState} from "react"
import type {ReactNode} from "react"

class SerializableSet<T> extends Set<T> {
  toJSON() {
    return Array.from(this);
  }
}

const filterDirections = [
    "ascending",
    "descending"
] as const;

type FilterDirections = typeof filterDirections[number];

const filterElements = [
  "date_created",
  "date_uploaded",
  "filename",
  "title",
] as const;

type FilterElement = typeof filterElements[number];

type Filter = {
    title?: string
    filter_element: FilterElement
    filter_direction: FilterDirections
    filter_tags: SerializableSet<string>
}

type FilterContextType = {
    filter: Filter,
    toggleTag: (tagId: string) => Set<string> 
    updateFilterParameter: <K extends keyof Filter>(param: K, val: Filter[K]) => void,    
    getDefaultFilter: () => Filter
    cycleFilterDirection: (direction: FilterDirections) => void
    cycleFilterElement: (element: FilterElement) => void
}

type FilterProviderProps = {
  children: ReactNode;
};

const FilterContext = createContext<FilterContextType | undefined>(undefined)

function getDefaultFilter(): Filter {
    return {
        filter_element: "date_created",
        filter_direction: "descending",
        filter_tags: new SerializableSet<string>()
    }
}

export function FilterProvider({children}: FilterProviderProps) {
    const [filter, setFilter] = useState<Filter>(getDefaultFilter())

    function updateFilterParameter<K extends keyof Filter>(param: K, val: Filter[K]) {
        setFilter(prev => ({ ...prev, [param]: val }))
    }

    //its a list so its small, i mean maybe if we get over 100 different filters :sob:
    function cycleFilterDirection(direction: FilterDirections) {
        const index = filterDirections.indexOf(direction);
        const newItem = filterDirections[(index + 1) % filterDirections.length]
        updateFilterParameter("filter_direction", newItem)
        return filterDirections[(index + 1) % filterDirections.length];
    }

    function cycleFilterElement(element: FilterElement) {
        const index = filterElements.indexOf(element);
        const newItem = filterElements[(index + 1) % filterElements.length]
        updateFilterParameter("filter_element", newItem)
        return newItem;    
    }

    function toggleTag(tagId: string): SerializableSet<string> {
        const newTags = new SerializableSet(filter.filter_tags); 

        if (newTags.has(tagId)) {
            newTags.delete(tagId); 
        } else {
            newTags.add(tagId); 
        }

        updateFilterParameter("filter_tags", newTags)
        return newTags;
    }

    return (
        <FilterContext.Provider
            value={{
                filter,
                toggleTag: toggleTag,
                updateFilterParameter: updateFilterParameter,
                getDefaultFilter: getDefaultFilter,
                cycleFilterDirection: cycleFilterDirection,
                cycleFilterElement: cycleFilterElement
            }}
        >
        {children}
        </FilterContext.Provider>
    )
}


export function useFilter() {
    const context = useContext(FilterContext)

    if (!context) {
        throw new Error("Filter Context not intialized!")
    }

    return context
}