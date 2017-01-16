<p align="center">
  <img width="220" src="https://raw.githubusercontent.com/cepave-f2e/owl-light/dev/assets/logo.png" />
</p>

<p align="center">
  <a href="https://circleci.com/gh/cepave-f2e/owl-light" alt="Build Status" target="_blank">
    <img src="https://img.shields.io/circleci/project/github/cepave-f2e/owl-light/dev.svg" />
  </a>
  <a href="https://codecov.io/gh/cepave-f2e/owl-light" alt="Coverage" target="_blank">
    <img src="https://img.shields.io/codecov/c/github/cepave-f2e/owl-light.svg" />
  </a>
  <a href="https://github.com/cepave-f2e/owl-light/releases">
    <img src="https://img.shields.io/github/tag/cepave-f2e/owl-light.svg" />
  </a>
  <img src="https://img.shields.io/github/license/cepave-f2e/owl-light.svg" />
</p>


# OWL Light

OWL Light is an [Open-Falcon](https://github.com/open-falcon) client-side project. It's based on [vue](https://github.com/vuejs/vue) and [vue-router](https://github.com/vuejs/vue-router) 2 worked as **SPA** (Single Page Application).

## Setup

Clone this project

```sh
git clone https://github.com/cepave-f2e/owl-light.git && cd owl-light
```

## Install dependencies

Recommend use `yarn` to install.

```sh
yarn install
```



## The Folder Structure

```
owl-light/
├── build/ (webpack build and dev server config)
├── config/ (all project configs)
├── src/ (source, core code base)
|    ├─── components/ (global common components)
|    ├─── containers/ (page container)
|    ├─── sass/ (global common sass/scss libraries)
|    ├─── store/ (Vuex store management)
|    └─── utils/ (global common utility functions)
├── .babelrc (babel complier config)
├── .editorconfig (editor config)
├── .eslintrc.js (eslint config)
├── .stylelintrc.js (stylelint config)
├── package.json
└── yarn.lock (modules cache file)
```


## Configure

All the related configs are in `/config`.


### API service

Configure you own API service, put the `OWL_LIGHT_API_BASE` environment variable


## OWL UI

[OWL UI](https://cepave-f2e.github.io/vue-owl-ui) is a **Component Design System** based on VueJS 2, used by Cepave to run in monitoring system and OWL Light



## Development

```sh
npm run dev
```

Open [http://localhost:8080](http://localhost:8080) to view it in the browser.


If you'd like to open browser automation, it just pass the `--open` arg.

```sh
npm run dev -- --open
```

## Build

The build files it'll output in `/dist`.

```sh
npm run build
```
