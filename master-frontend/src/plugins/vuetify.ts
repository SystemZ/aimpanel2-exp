import Vue from 'vue';
import Vuetify from 'vuetify/lib';
import colors from 'vuetify/lib/util/colors';

Vue.use(Vuetify);

export default new Vuetify({
    theme: {
        // select dark theme by default
        dark: true,
        options: {
            // enabling this can add some KBs to page size
            customProperties: false,
            minifyTheme(css) {
                return process.env.NODE_ENV === 'production'
                    ? css.replace(/[\r\n|\r|\n]/g, '')
                    : css;
            },
        },
        themes: {
            light: {
                primary: '#1976D2',
                secondary: '#424242',
                accent: '#82B1FF',
                error: '#FF5252',
                info: '#2196F3',
                success: '#4CAF50',
                warning: '#FFC107',
            },
            dark: {
                // custom
                toolbar: colors.grey.darken4,
                background: colors.grey.darken4,
                cardBg: colors.grey.darken3,
                buttons: colors.green.base,
                buttonsDestroy: colors.red.base,
                toolbarLoader: colors.green.base,
                // original
                anchor: colors.red.base,
            },
        },
    },
    icons: {
        iconfont: 'fa4',
    },
});
