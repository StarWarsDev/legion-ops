import React, { createContext, useReducer, useContext } from 'react'
import { User } from './model/app';

interface Props {
    children: React.ReactChildren | Element | React.ReactChild;
}

interface IState {
    user: User;
}

interface IAppContext {
    state: typeof initialState;
    dispatch: (action: Action) => void;
}

type Action =
    | {type: "authenticated", user: User};

const initialState = {
    user: { authenticated: false, username: "", name: "", picture: "" }
};

const AppContext = createContext<IAppContext>({
    state: initialState,
    dispatch: () => {}
});

const reducer = (state: IState, action: Action) => {
    switch (action.type) {
        case "authenticated":
            return {
                ...state,
                user: action.user
            };
        default:
            return state
    }
};

export const AppDataProvider = ({ children }: Props) => {
    const [state, dispatch] = useReducer(reducer, initialState);

    return (
        <AppContext.Provider value={{ state, dispatch }}>
            {children}
        </AppContext.Provider>
    )
}

export const useAppStore = () => useContext(AppContext);