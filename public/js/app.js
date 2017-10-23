(function() {
    var _ = {
        forEach: function (array, callback, scope) {
            for (var i = 0; i < array.length; i++) {
                callback.call(scope, i, array[i]);
            }
        }
    };


    // Top menu dropdowns
    _.forEach(document.querySelectorAll('.navbar-item.has-dropdown'), function(i, container) {
        var button = container.querySelector('.navbar-link');
        button.addEventListener('click', function() {
            container.classList.toggle('is-active');
        });
    });

    // Hamburger menu for mobile devices
    _.forEach(document.querySelectorAll('.navbar-burger'), function(button) {
        var menu = document.querySelector('#' + (button.getAttribute ? button.getAttribute('data-target') : button['data-target']));
        button.addEventListener('click', function() {
            button.classList.toggle('is-active');
            menu.classList.toggle('is-active');
        });
    });

    var closeDropdownsIfOutside = function (target) {
        _.forEach(document.querySelectorAll('.has-dropdown.is-active'), function(i, container) {
            if (target === null || !container.contains(target)) {
                container.classList.toggle('is-active');
            }
        });
    };

    // Close all open dropdowns if clicked outside
    document.querySelector('body, .site').addEventListener('click', function(event) {
        closeDropdownsIfOutside(event.target)
    });

    // Close all open dropdowns on escape button
    document.querySelector('body, .site').addEventListener('keyup', function(event) {
        if (event.keyCode === 27 || event.key === 'Escape' || event.key == 'Esc') {
            closeDropdownsIfOutside(null);
        }
    });
})();
