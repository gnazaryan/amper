export const parseBoolean = (value) => {
    if (value != null && (value === true || value === 'true' || value === 1 || value === '1' || value === 'yes')) {
        return true;
    }
    return false;
}