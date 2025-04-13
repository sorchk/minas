import { Result } from "@/api/ajax";
import { Ref, ref } from "vue"
import { FormItemRule } from "naive-ui";
import { t } from "@/locales";

export function useForm<T>(form: Ref, action: () => Promise<Result<T>>, success?: (data: T) => void) {
    const submiting = ref(false)
    async function submit(e: Event) {
        e.preventDefault();
        form.value.validate(async (errors: any) => {
            if (errors) {
                return
            }
            submiting.value = true;
            try {
                let r = await action()
                success ? success(<T>r.data) : window.message.info(t('texts.action_success'));
            } finally {
                submiting.value = false;
            }
        });
    }

    return { submit, submiting }
}

export function requiredRule(field?: string, message?: string): FormItemRule {
    return {
        required: true,
        message: formatMessage(field, message ?? t('tips.required_rule')),
        trigger: ["input", "blur", 'change'],
    }
}
export function requiredNumberRule(field?: string, message?: string): FormItemRule {
    return {
        type: 'number',
        required: true,
        message: formatMessage(field, message ?? t('tips.required_rule')),
        trigger: ["input", "blur", 'change'],
    }
}

export function customRule(validator: (rule: any, value: any) => boolean, message?: string, field?: string, required?: boolean): FormItemRule {
    return createRule(validator, message, field, required)
}

export function portValidator(x: number) {
    return x >= 1 && x <= 65535;
}
export function portRule(field?: string): FormItemRule  {
    return {
        message: formatMessage(field, t('tips.port_rule')),
        trigger: ["input", "blur"],
        validator(rule: any, value: string): boolean {
            if (!value) {
                return true;
            }
            try{
                const x = parseInt(value);
                return x >= 1 && x <= 65535;
            }catch(e){
                return false;
            }
        },
    }
}
export function ipDomainRule(field?: string): FormItemRule {
    const domainRegex = /^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$/i;
    const ipRegex = /(^((2[0-4]\d.)|(25[0-5].)|(1\d{2}.)|(\d{1,2}.))((2[0-5]{2}.)|(1\d{2}.)|(\d{1,2}.){2})((1\d{2})|(2[0-5]{2})|(\d{1,2})))/;
    return {
        message: formatMessage(field, t('tips.ip_domain_rule')),
        trigger: ["input", "blur"],
        validator(rule: any, value: string): boolean {
            return !value || ipRegex.test(value)|| domainRegex.test(value)
        },
    };
}
export function ipRule(field?: string): FormItemRule {
    const ipRegex = /(^((2[0-4]\d.)|(25[0-5].)|(1\d{2}.)|(\d{1,2}.))((2[0-5]{2}.)|(1\d{2}.)|(\d{1,2}.){2})((1\d{2})|(2[0-5]{2})|(\d{1,2})))/;
    return regexRule(ipRegex, t('tips.ip_rule'), field)
}

export function emailRule(field?: string): FormItemRule {
    const reg = /^([a-zA-Z0-9]+[-_\.]?)+@[a-zA-Z0-9]+\.[a-z]+$/;
    return regexRule(reg, t('tips.email_rule'), field)
}

export function phoneRule(field?: string): FormItemRule {
    const reg = /^[1][3,4,5,7,8][0-9]{9}$/;
    return regexRule(reg, t('tips.phone_rule'), field)
}


export function lengthRule(min: number, max: number, field?: string): FormItemRule {
    return createRule((rule: any, value: string): boolean => {
        return value.length >= min && value.length <= max
    }, t('tips.length_rule', { min, max }), field)
}

export function passwordRule(field?: string): FormItemRule {
    const reg = /^[\w\W]+$/;
    return regexRule(reg, t('tips.password_rule'), field)
}

export function regexRule(reg: RegExp, message?: string, field?: string): FormItemRule {
    return {
        message: formatMessage(field, message),
        trigger: ["input", "blur"],
        validator(rule: any, value: string): boolean {
            return !value || reg.test(value)
        },
    };
}

function createRule(validator: (rule: any, value: string) => boolean, message?: string, field?: string, required?: boolean): FormItemRule {
    return {
        required: required,
        message: formatMessage(field, message),
        trigger: ["input", "blur"],
        validator,
    };
}

function formatMessage(field?: string, message?: string) {
    return field ? `${field}: ${message}` : message
}